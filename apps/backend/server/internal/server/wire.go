package server

import (
	"net/http"
	"qingmo-server/internal/config"
	"qingmo-server/internal/controller"
	"qingmo-server/internal/repository"
	"qingmo-server/internal/security"
	"qingmo-server/internal/service"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func wireAuth(
	db *gorm.DB,
	tokenManager *security.TokenManager,
	redis *redis.Client,
) *controller.AuthController {
	users := repository.NewUserPostgresRepository(db)
	refreshSession := repository.NewRefreshSessionRedisRepository(redis)

	authService := service.NewAuthService(users, tokenManager, refreshSession)

	return controller.NewAuthController(
		authService,
		tokenManager,
		refreshSession,
	)
}

func wireTokenManager(
	cfg config.JWTConfig,
	redis *redis.Client,
) *security.TokenManager {
	keyStore := repository.NewSignKeyRedisRepository(redis)

	return security.NewTokenManager(
		keyStore,
		cfg.Issuer,
		cfg.Audience,
		cfg.AccessTTL,
		cfg.RefreshTTL,
	)
}

func wireUser(db *gorm.DB) *controller.UserController {
	users := repository.NewUserPostgresRepository(db)
	userService := service.NewUserService(users)
	return controller.NewUserController(
		userService,
	)
}

func wireCollaboration(
	db *gorm.DB,
	redis *redis.Client,
	wsHandler http.Handler,
) *controller.CollaborationController {
	documents := repository.NewDocumentAccessPostgresRepository(db)
	tickets := repository.NewCollaborationTicketRedisRepository(redis)
	collabService := service.NewCollaborationService(documents, tickets)
	return controller.NewCollaborationController(
		collabService,
		wsHandler,
	)
}
