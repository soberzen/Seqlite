package controller

import (
	"qingmo-server/internal/dto"
	"qingmo-server/internal/errs"
	"qingmo-server/internal/response"
	"qingmo-server/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SpaceController struct {
	spaces *service.SpaceService
}

func NewSpaceController(spaces *service.SpaceService) *SpaceController {
	return &SpaceController{spaces: spaces}
}

func (s *SpaceController) Create(c *gin.Context, userID uint) {
	var req dto.CreateSpaceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BindError(c, err)
		return
	}

	space, err := s.spaces.Create(c.Request.Context(), userID, req)
	if err != nil {
		response.WriteError(c, err)
		return
	}
	response.Success(c, space)
}

func (s *SpaceController) List(c *gin.Context, userID uint) {
	spaces, err := s.spaces.ListVisibleSpaces(c.Request.Context(), userID)
	if err != nil {
		response.WriteError(c, err)
		return
	}
	response.Success(c, spaces)
}

func (s *SpaceController) ListTrashed(c *gin.Context, userID uint) {
	spaces, err := s.spaces.ListTrashedSpaces(c.Request.Context(), userID)
	if err != nil {
		response.WriteError(c, err)
		return
	}
	response.Success(c, spaces)
}

func (s *SpaceController) Get(c *gin.Context, userID uint) {
	publicID, ok := spacePublicID(c)
	if !ok {
		return
	}
	space, err := s.spaces.GetSpace(c.Request.Context(), userID, publicID)
	if err != nil {
		response.WriteError(c, err)
		return
	}
	response.Success(c, space)
}

func (s *SpaceController) Update(c *gin.Context, userID uint) {
	publicID, ok := spacePublicID(c)
	if !ok {
		return
	}
	var req dto.UpdateSpaceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BindError(c, err)
		return
	}
	space, err := s.spaces.UpdateSpaceInfo(
		c.Request.Context(),
		userID,
		publicID,
		req,
	)
	if err != nil {
		response.WriteError(c, err)
		return
	}
	response.Success(c, space)
}

func (s *SpaceController) ListDocuments(c *gin.Context, userID uint) {
	publicID, ok := spacePublicID(c)
	if !ok {
		return
	}
	var req dto.ListSpaceDocumentsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BindError(c, err)
		return
	}
	page, err := s.spaces.ListDocuments(
		c.Request.Context(),
		userID,
		publicID,
		req,
	)
	if err != nil {
		response.WriteError(c, err)
		return
	}
	response.Success(c, page)
}

func spacePublicID(c *gin.Context) (uuid.UUID, bool) {
	id, err := uuid.Parse(c.Param("spaceId"))
	if err != nil || id == uuid.Nil {
		response.WriteError(c, errs.New(errs.KindInvalid, "invalid space id"))
		return uuid.Nil, false
	}
	return id, true
}
