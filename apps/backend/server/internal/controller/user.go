package controller

import (
	"qingmo-server/internal/dto"
	"qingmo-server/internal/errs"
	"qingmo-server/internal/middleware"
	"qingmo-server/internal/response"
	"qingmo-server/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService *service.UserService
}

func NewUserController(
	userService *service.UserService,
) *UserController {
	return &UserController{
		userService: userService,
	}
}

func (u *UserController) Profile(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.WriteError(c, errs.New(errs.KindUnauthorized, "invalid token"))
		return
	}
	result, err := u.userService.Profile(c, userID)
	if err != nil {
		response.WriteError(c, err)
		return
	}
	response.Success(c, result)
}

func (u *UserController) UpdateProfile(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.WriteError(c, errs.New(errs.KindUnauthorized, "invalid token"))
		return
	}

	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BindError(c, err)
		return
	}
	result, err := u.userService.UpdateProfile(c.Request.Context(), userID, req)
	if err != nil {
		response.WriteError(c, err)
		return
	}
	response.Success(c, result)
}

func (u *UserController) PublicProfile(c *gin.Context) {
	// A PostgreSQL BIGINT has 63 bits available for positive IDs.
	id, err := strconv.ParseUint(c.Param("id"), 10, min(strconv.IntSize, 63))
	if err != nil || id == 0 {
		response.WriteError(c, errs.New(errs.KindInvalid, "invalid user id"))
		return
	}

	result, err := u.userService.PublicProfile(c.Request.Context(), uint(id))
	if err != nil {
		response.WriteError(c, err)
		return
	}
	response.Success(c, result)
}
