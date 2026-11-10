package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"qingmo-server/internal/config"
	"qingmo-server/internal/server"
)

func main() {
	if err := config.LoadEnv("./"); err != nil {
		// 加载环境变量失败
		fmt.Fprintf(os.Stderr, "failed to load env: %v\n", err)
		os.Exit(-1)
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		// 加载配置失败
		fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
		os.Exit(-1)
	}
	// 加载配置成功
	// 初始化 logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	server, err := server.NewServer(cfg, logger)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create server: %v\n", err)
		os.Exit(-1)
	}

	done := make(chan bool, 1)
	go gracefulShutdown(server, done)

	err = server.ListenAndServe()

	if err != nil && err != http.ErrServerClosed {
		panic(fmt.Sprintf("http server error: %s", err))
	}

	<-done // 等待服务器关闭
	log.Println("Graceful shutdown completed")
}

func gracefulShutdown(httpServer *http.Server, done chan bool) {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	<-ctx.Done() // 等待信号

	log.Println("shutting down gracefully, press Ctrl+C again to force")
	stop()
	// 延迟5秒，确保所有请求完成关闭
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Println("error shutting down server server:", err)
	}
	log.Println("Server exiting...")

	done <- true
}
