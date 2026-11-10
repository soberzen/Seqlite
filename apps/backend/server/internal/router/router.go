package router

import (
	"log/slog"
	"net/http"

	"qingmo-server/internal/controller"
	"qingmo-server/internal/middleware"

	"github.com/gin-gonic/gin"
)

func NewRoute(
	authController *controller.AuthController,
	logger *slog.Logger,
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

	registerAuthRoutes(public, authController)
	registerProtectedRoutes(protected, authController)

	return router
}
