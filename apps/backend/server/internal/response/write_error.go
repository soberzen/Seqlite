package response

import (
	"context"
	"errors"
	"qingmo-server/internal/errs"

	"github.com/gin-gonic/gin"
)

func WriteError(c *gin.Context, err error) {
	if err == nil {
		return
	}
	// 客户端已断开：不写响应
	if errors.Is(err, context.Canceled) ||
		errors.Is(err, context.DeadlineExceeded) {
		c.Abort()
		return
	}

	kind := errs.KindOf(err)
	httpStatus := kind.HTTPStatus()

	// 500 内部服务器错误
	if kind == errs.KindInternal {
		c.Error(err)
		AbortWithError(c, httpStatus, CodeInternal, "internal server error")
		return
	}
	// 4xx：暴露消息 + 业务码
	AbortWithError(c, httpStatus, businessCode(err), err.Error())
}

func businessCode(err error) int {
	switch errs.KindOf(err) {
	case errs.KindInvalid:
		return CodeInvalidParam
	case errs.KindUnauthorized:
		return CodeUnauthorized
	case errs.KindForbidden:
		return CodeForbidden
	case errs.KindNotFound:
		return CodeNotFound
	case errs.KindConflict:
		return CodeConflict
	default:
		return CodeInternal
	}
}
