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
	authController *controller.AuthController,
	logger *slog.Logger,
	tokenManager *security.TokenManager,
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

	return router
}
