package router

import (
	"log/slog"
	"net/http"

	"qingmo-server/internal/controller"
	"qingmo-server/internal/errs"
	"qingmo-server/internal/middleware"
	"qingmo-server/internal/response"
	"qingmo-server/internal/security"

	"github.com/gin-gonic/gin"
)

func NewRoute(
	logger *slog.Logger,
	tokenManager *security.TokenManager,
	authController *controller.AuthController,
	userController *controller.UserController,
	spaceController *controller.SpaceController,
	collabationController *controller.CollaborationController,
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
	api := router.Group("/v1/api")

	public := api.Group("")
	protected := api.Group("")
	protected.Use(middleware.Auth(tokenManager))

	registerAuthRoutes(public, authController)
	registerProtectedRoutes(protected, authController)
	registerUserRoutes(protected, userController)
	registerSpaceRoutes(protected, spaceController)

	// websocket
	registerCollaborationRoutes(public, protected, collabationController)

	return router
}

func withAuthenticatedUser(handler func(*gin.Context, uint)) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := middleware.GetUserID(c)
		if !ok {
			response.WriteError(
				c,
				errs.New(errs.KindUnauthorized, "invalid token"),
			)
			return
		}
		handler(c, userID)
	}
}
