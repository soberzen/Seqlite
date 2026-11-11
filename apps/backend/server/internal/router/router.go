package router

import (
	"log/slog"
	"net/http"

	"qingmo-server/internal/controller"
	"qingmo-server/internal/middleware"
	"qingmo-server/internal/security"

	"github.com/gin-gonic/gin"
)

func NewRoute(
	logger *slog.Logger,
	tokenManager *security.TokenManager,
	authController *controller.AuthController,
	userController *controller.UserController,
) http.Handler {
	router := gin.New()
	// Middleware
	registerMiddleware(
		router,
		gin.Recovery(),
		middleware.RequestID(),
		middleware.Logger(logger),
	)

	// API routes
	api := router.Group("/api")

	public := api.Group("")
	protected := api.Group("")
	protected.Use(middleware.Auth(tokenManager))

	registerAuthRoutes(public, authController)
	registerProtectedRoutes(protected, authController)
	registerUserRoutes(protected, userController)

	return router
}
