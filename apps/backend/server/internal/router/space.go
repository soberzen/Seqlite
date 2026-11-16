package router

import (
	"qingmo-server/internal/controller"

	"github.com/gin-gonic/gin"
)

func registerSpaceRoutes(
	group *gin.RouterGroup,
	spaces *controller.SpaceController,
) {
	routes := group.Group("/space")
	routes.POST("", withAuthenticatedUser(spaces.Create))
	routes.GET("", withAuthenticatedUser(spaces.List))
	routes.GET("/trash", withAuthenticatedUser(spaces.ListTrashed))
	routes.GET("/:spaceId", withAuthenticatedUser(spaces.Get))
	routes.PATCH("/:spaceId", withAuthenticatedUser(spaces.Update))
	routes.GET(
		"/:spaceId/documents",
		withAuthenticatedUser(spaces.ListDocuments),
	)
}
