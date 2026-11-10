package router

import "github.com/gin-gonic/gin"

func registerMiddleware(router *gin.Engine, middlewares ...gin.HandlerFunc) {
	for _, middleware := range middlewares {
		router.Use(middleware)
	}
}
