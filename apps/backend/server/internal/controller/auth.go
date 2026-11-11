package controller

import (
	"qingmo-server/internal/dto"
	"qingmo-server/internal/errs"
	"qingmo-server/internal/repository"
	"qingmo-server/internal/response"
	"qingmo-server/internal/security"
	"qingmo-server/internal/service"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService    *service.AuthService
	tokenManager   *security.TokenManager
	refreshSession *repository.RefreshSessionRedisRepository
}

const (
	refreshCookieName = "refresh_token"
	refreshCookiePath = "/api/auth" // 和路由保持一致
)

func NewAuthController(
	authService *service.AuthService,
	tokenManager *security.TokenManager,
	refreshSessionRedisRepository *repository.RefreshSessionRedisRepository,
) *AuthController {
	return &AuthController{
		authService:    authService,
		tokenManager:   tokenManager,
		refreshSession: refreshSessionRedisRepository,
	}
}

func (a *AuthController) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BindError(c, err)
		return
	}
	result, err := a.authService.Login(c, &req)
	if err != nil {
		response.WriteError(c, err)
		return
	}
	c.SetCookie(
		refreshCookieName,
		result.RefreshToken,
		int(a.tokenManager.RefreshTTL.Seconds()),
		refreshCookiePath,
		"",
		false, // 开发环境不开启 HTTPS
		true,
	)

	response.Success(c, dto.LoginResponse{
		AccessToken: result.AccessToken,
	})

}
func (a *AuthController) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BindError(c, err)
		return
	}
	result, err := a.authService.Register(c, req)
	if err != nil {
		response.WriteError(c, err)
		return
	}
	response.Success(c, dto.RegisterResponse{
		UserName:  result.UserName,
		Email:     result.Email,
		AvatarURL: result.AvatarURL,
	})
}
func (a *AuthController) Refresh(c *gin.Context) {
	refreshToken, err := c.Cookie(refreshCookieName)
	if err != nil {
		response.WriteError(
			c,
			errs.Wrap(errs.KindUnauthorized, "missing refresh token", err),
		)
		return
	}

	result, err := a.authService.Refresh(
		c,
		refreshToken,
	)
	if err != nil {
		response.WriteError(c, err)
		return
	}

	c.SetCookie(
		refreshCookieName,
		result.RefreshToken,
		int(a.tokenManager.RefreshTTL.Seconds()),
		refreshCookiePath,
		"",
		false, // 开发环境不开启 HTTPS
		true,
	)

	response.Success(c, dto.LoginResponse{
		AccessToken: result.AccessToken,
	})
}

func (a *AuthController) Logout(c *gin.Context) {
	refreshToken, err := c.Cookie(refreshCookieName)
	c.SetCookie(
		refreshCookieName,
		"",
		-1,
		refreshCookiePath,
		"",
		false,
		true,
	)
	if err != nil {
		response.Success[any](
			c,
			nil,
		)
		return
	}
	if err := a.authService.Logout(c, refreshToken); err != nil {
		response.WriteError(c, err)
		return
	}
	response.Success[any](
		c,
		nil,
	)
}
