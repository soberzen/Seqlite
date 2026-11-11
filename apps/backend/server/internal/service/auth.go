package service

import (
	"context"
	"qingmo-server/internal/dto"
	"qingmo-server/internal/errs"
	"qingmo-server/internal/model"
	"qingmo-server/internal/repository"
	"qingmo-server/internal/security"
	"strconv"
	"strings"
)

type AuthService struct {
	users          *repository.UserPostgresRepository
	tokenManager   *security.TokenManager
	refreshSession *repository.RefreshSessionRedisRepository
}

func NewAuthService(
	users *repository.UserPostgresRepository,
	tokenManager *security.TokenManager,
	refreshSessionRedisRepository *repository.RefreshSessionRedisRepository,
) *AuthService {
	return &AuthService{
		users:          users,
		tokenManager:   tokenManager,
		refreshSession: refreshSessionRedisRepository,
	}
}

type LoginResult struct {
	AccessToken  string
	RefreshToken string
	SessionID    string
}

func (a *AuthService) Login(
	ctx context.Context,
	req *dto.LoginRequest,
) (*LoginResult, error) {
	user, err := a.users.GetUserByEmail(ctx, strings.ToLower(req.Email))

	if err != nil {
		return nil, err
	}
	valid, err := security.VerifyPassword(req.Password, user.Password)
	if err != nil {
		return nil, err
	}
	if !valid {
		return nil, errs.New(
			errs.KindUnauthorized,
			"invalid email or password",
		)
	}
	// 登录成功，生成token
	sessionId, err := security.GenerateSessionID()
	if err != nil {
		return nil, errs.Wrap(
			errs.KindInternal,
			"fail generate session ID",
			err,
		)
	}

	accessToken, _, err := a.tokenManager.SignJWT(
		ctx,
		uint(user.ID),
		security.TokenTypeAccess,
		sessionId,
	)
	if err != nil {
		return nil, err
	}
	refreshToken, claims, err := a.tokenManager.SignJWT(
		ctx,
		uint(user.ID),
		security.TokenTypeRefresh,
		sessionId,
	)
	if err != nil {
		return nil, err
	}
	// 创建refresh session
	if err := a.refreshSession.Create(
		ctx,
		sessionId,
		&repository.RefreshSession{
			UserID:    uint(user.ID),
			JTI:       claims.ID,
			IssuedAt:  claims.IssuedAt.Time,
			ExpiresAt: claims.ExpiresAt.Time,
		},
	); err != nil {
		return nil, err
	}
	return &LoginResult{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		SessionID:    sessionId,
	}, nil
}

func (a *AuthService) Register(
	ctx context.Context,
	req dto.RegisterRequest,
) (*dto.RegisterResponse, error) {
	exists, err := a.users.ExistsByEmail(ctx, strings.ToLower(req.Email))
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errs.New(errs.KindConflict, "email is already taken")
	}
	encryptedPassword, err := security.EncryptPassword(req.Password)
	if err != nil {
		return nil, errs.Wrap(errs.KindInternal, "fail encrypt password", err)
	}
	user, err := a.users.Create(ctx, &model.Users{
		UserName:  req.UserName,
		Email:     strings.ToLower(req.Email),
		Password:  encryptedPassword,
		AvatarURL: &req.AvatarURL,
	})
	if err != nil {
		return nil, err
	}
	return &dto.RegisterResponse{
		UserName:  user.UserName,
		Email:     user.Email,
		AvatarURL: *user.AvatarURL,
	}, nil
}

func (a *AuthService) Logout(ctx context.Context, refreshToken string) error {
	claims, err := a.tokenManager.VerifyJWT(
		ctx,
		refreshToken,
		security.TokenTypeRefresh,
	)
	if err != nil {
		return err
	}
	if err := a.refreshSession.Delete(ctx, claims.SessionID); err != nil {
		return errs.Wrap(errs.KindInternal, "fail logout", err)
	}
	return nil
}

func (a *AuthService) Refresh(
	ctx context.Context,
	refreshToken string,
) (*LoginResult, error) {
	claims, err := a.tokenManager.VerifyJWT(
		ctx,
		refreshToken,
		security.TokenTypeRefresh,
	)
	if err != nil {
		return nil, errs.Wrap(
			errs.KindUnauthorized,
			"invalid refresh token",
			err,
		)
	}
	oldSessionID := claims.SessionID
	expectedOldJTI := claims.ID
	userID, err := strconv.ParseUint(claims.Subject, 10, 64)
	if err != nil {
		return nil, errs.Wrap(
			errs.KindUnauthorized,
			"invalid subject in token",
			err,
		)
	}
	if userID == 0 {
		return nil, errs.New(errs.KindUnauthorized, "invalid token")
	}

	newSessionID, err := security.GenerateSessionID()
	if err != nil {
		return nil, errs.Wrap(
			errs.KindInternal,
			"fail generate session ID",
			err,
		)
	}
	accessToken, _, err := a.tokenManager.SignJWT(
		ctx,
		uint(userID),
		security.TokenTypeAccess,
		newSessionID,
	)
	if err != nil {
		return nil, err
	}
	refreshToken, newRefreshClaims, err := a.tokenManager.SignJWT(
		ctx,
		uint(userID),
		security.TokenTypeRefresh,
		newSessionID,
	)
	if err != nil {
		return nil, err
	}

	newRefreshSession := &repository.RefreshSession{
		UserID:    uint(userID),
		JTI:       newRefreshClaims.ID,
		IssuedAt:  newRefreshClaims.IssuedAt.Time,
		ExpiresAt: newRefreshClaims.ExpiresAt.Time,
	}

	if err := a.refreshSession.Rotate(
		ctx,
		oldSessionID,
		expectedOldJTI,
		newSessionID,
		newRefreshSession,
	); err != nil {
		if errs.IsKind(err, errs.KindNotFound) ||
			errs.IsKind(err, errs.KindConflict) {
			return nil, errs.Wrap(
				errs.KindUnauthorized,
				"refresh session invalid",
				err,
			)
		}
		return nil, err
	}

	return &LoginResult{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		SessionID:    newSessionID,
	}, nil
}
