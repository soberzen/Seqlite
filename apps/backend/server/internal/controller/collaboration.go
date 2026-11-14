package controller

import (
	"net/http"
	"qingmo-server/internal/errs"
	"qingmo-server/internal/middleware"
	"qingmo-server/internal/realtime"
	"qingmo-server/internal/response"
	"qingmo-server/internal/service"

	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
)

type CollaborationController struct {
	service   *service.CollaborationService
	wsHandler http.Handler
}

func NewCollaborationController(
	svc *service.CollaborationService,
	wsHandler http.Handler,
) *CollaborationController {
	return &CollaborationController{
		service:   svc,
		wsHandler: wsHandler,
	}
}

func collaborationDocumentID(c *gin.Context) (uuid.UUID, error) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil || id == uuid.Nil {
		return uuid.Nil, errs.New(errs.KindInvalid, "invalid document")
	}
	return id, nil
}

func (h *CollaborationController) IssueTicket(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok || int64(userID) <= 0 {
		response.WriteError(c, errs.New(errs.KindUnauthorized, "invalid token"))
		return
	}
	publicID, err := collaborationDocumentID(c)
	if err != nil {
		response.WriteError(c, err)
		return
	}
	result, err := h.service.IssueTicket(c, publicID, userID)
	if err != nil {
		response.WriteError(c, err)
		return
	}

	c.Header("Cache-Control", "no-store")
	response.Success[any](c, result)
}

func (h *CollaborationController) Connect(c *gin.Context) {
	publicID, err := collaborationDocumentID(c)
	if err != nil {
		response.WriteError(c, err)
		return
	}
	ticket := c.Query("ticket")
	if ticket == "" || len(ticket) > 128 {
		response.WriteError(
			c,
			errs.New(errs.KindUnauthorized, "invalid websocket ticket"),
		)
		return
	}
	access, err := h.service.AuthorizeConnect(c, publicID, ticket)
	if err != nil {
		response.WriteError(c, err)
		return
	}
	req := realtime.WithDocumentAccess(
		c.Request,
		access.DocumentID,
		access.CanEdit,
	)
	h.wsHandler.ServeHTTP(c.Writer, req)
}
