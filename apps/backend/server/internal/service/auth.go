package service

import (
	"context"
	"qingmo-server/internal/dto"
	"qingmo-server/internal/errs"
	"qingmo-server/internal/model"
	"qingmo-server/internal/repository"
	"qingmo-server/internal/security"
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
