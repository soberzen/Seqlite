package router

import (
	"qingmo-server/internal/controller"

	"github.com/gin-gonic/gin"
)

func registerUserRoutes(
	group *gin.RouterGroup,
	userController *controller.UserController,
) {
	routes := group.Group("/user")
	routes.GET("/me", userController.Profile)
	routes.PATCH("/me", userController.UpdateProfile)
	group.GET("/users/:id", userController.PublicProfile)
}
