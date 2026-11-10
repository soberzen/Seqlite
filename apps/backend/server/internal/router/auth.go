package router

import (
	"qingmo-server/internal/controller"

	"github.com/gin-gonic/gin"
)

func registerAuthRoutes(
	group *gin.RouterGroup,
	authController *controller.AuthController,
) {
	routes := group.Group("/auth")
	routes.POST("/login", authController.Login)
	routes.POST("/register", authController.Register)
	routes.POST("/refresh", authController.Refresh)
}

func registerProtectedRoutes(
	group *gin.RouterGroup,
	authController *controller.AuthController,
) {
	routes := group.Group("/auth")
	routes.POST("/logout", authController.Logout)
}
