package controller

import (
	"qingmo-server/internal/errs"
	"qingmo-server/internal/middleware"
	"qingmo-server/internal/response"
	"qingmo-server/internal/service"

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
