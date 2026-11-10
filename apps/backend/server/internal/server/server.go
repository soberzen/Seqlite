package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"qingmo-server/internal/config"
	"qingmo-server/internal/database"
	"qingmo-server/internal/router"
)

func NewServer(cfg *config.Config, logger *slog.Logger) (*http.Server, error) {
	postgresDB := database.NewPostgres(cfg.Postgres)
	redisDB := database.NewRedis(cfg.Redis)

	if postgresDB == nil {
		return nil, fmt.Errorf("can not connect postgres")
	}
	if redisDB == nil {
		return nil, fmt.Errorf("can not connect redis")
	}
	tokenManager := wireTokenManager(cfg.JWT, redisDB.DB())

	authController := wireAuth(postgresDB.DB(), tokenManager, redisDB.DB())

	handler := router.NewRoute(authController, logger)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.HTTP.Port),
		Handler:      handler,
		IdleTimeout:  time.Minute,      // 1分钟内没有请求，关闭连接
		ReadTimeout:  10 * time.Second, // 10秒内没有读取到请求，关闭连接
		WriteTimeout: 30 * time.Second, // 30秒内没有写入响应，关闭连接
	}

	return server, nil
}
