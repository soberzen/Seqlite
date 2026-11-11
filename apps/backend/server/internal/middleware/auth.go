package middleware

import (
	"qingmo-server/internal/errs"
	"qingmo-server/internal/response"
	"qingmo-server/internal/security"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	ctxKeyUserID    = "user_id"
	ctxKeySessionID = "session_id"
	bearerPrefix    = "Bearer "
)

func Auth(tokenManager *security.TokenManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头中获取token
		tokenHeader := c.GetHeader("Authorization")
		if tokenHeader == "" {
			response.WriteError(
				c,
				errs.New(
					errs.KindUnauthorized,
					"authorization header is required",
				),
			)
			return
		}
		if !strings.HasPrefix(tokenHeader, bearerPrefix) {
			response.WriteError(c, errs.New(
				errs.KindUnauthorized,
				"authorization header must use Bearer scheme",
			))
			return
		}

		token := strings.TrimSpace(
			strings.TrimPrefix(tokenHeader, bearerPrefix),
		)
		if token == "" {
			response.WriteError(c, errs.New(
				errs.KindUnauthorized,
				"bearer token is empty",
			))
			return
		}
		// 验证token
		claims, err := tokenManager.VerifyJWT(
			c.Request.Context(),
			token,
			security.TokenTypeAccess,
		)
		if err != nil {
			response.WriteError(c, err)
			return
		}

		userID, err := strconv.ParseUint(claims.Subject, 10, 64)
		if err != nil {
			response.WriteError(c, errs.Wrap(
				errs.KindUnauthorized,
				"invalid subject in token",
				err,
			))
			return
		}
		if userID == 0 {
			response.WriteError(
				c,
				errs.New(errs.KindUnauthorized, "invalid token"),
			)
			return
		}

		c.Set(ctxKeyUserID, uint(userID))
		c.Set(ctxKeySessionID, claims.SessionID)

		c.Next()
	}
}

func GetUserID(c *gin.Context) (uint, bool) {
	val, exists := c.Get(ctxKeyUserID)
	if !exists {
		return 0, false
	}
	id, ok := val.(uint)
	return id, ok
}
