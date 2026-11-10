package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Result[T any] struct {
	Code      int    `json:"code"`
	Message   string `json:"message"`
	Data      T      `json:"data,omitempty"`
	RequestID string `json:"request_id,omitempty"`
}

func Success[T any](c *gin.Context, data T) {
	c.JSON(http.StatusOK, Result[T]{
		Code:      http.StatusOK,
		Message:   "success",
		Data:      data,
		RequestID: c.GetString("request_id"),
	})
}

func Error(c *gin.Context, httpStatus, code int, message string) {
	c.JSON(httpStatus, Result[any]{
		Code:      code,
		Message:   message,
		Data:      nil,
		RequestID: c.GetString("request_id"),
	})
}

func ErrorWithData(
	c *gin.Context,
	httpStatus,
	code int,
	message string,
	data any,
) {
	c.AbortWithStatusJSON(httpStatus, Result[any]{
		Code:      code,
		Message:   message,
		Data:      data,
		RequestID: c.GetString("request_id"),
	})
}

func AbortWithError(c *gin.Context, httpStatus, code int, message string) {
	Error(c, httpStatus, code, message)
	c.Abort()
}
