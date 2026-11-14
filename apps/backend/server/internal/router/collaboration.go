package router

import (
	"qingmo-server/internal/controller"

	"github.com/gin-gonic/gin"
)

func registerCollaborationRoutes(
	api,
	protect *gin.RouterGroup,
	collab *controller.CollaborationController,
) {
	const path = "/ws/document/:id"

	protect.POST(path+"/ticket", collab.IssueTicket)
	api.GET(path, collab.Connect)
}
