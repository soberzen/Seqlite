package database

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"time"

	"qingmo-server/internal/config"

	"github.com/redis/go-redis/v9"
)

type RedisService interface {
	Close() error
	DB() *redis.Client
}

type redisService struct {
	db *redis.Client
}

var redisInstance *redisService

func NewRedis(cfg config.RedisConfig) RedisService {
	if redisInstance != nil {
		return redisInstance
	}

	options := &redis.Options{
		Addr:                  fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password:              cfg.Password,
		DB:                    cfg.DB,
		ContextTimeoutEnabled: true,
	}
	if cfg.TLSEnabled {
		options.TLSConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
			ServerName: cfg.Host,
		}
	}
	client := redis.NewClient(options)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		log.Fatalf("can not connect redis: %v", err)
		return nil
	}
	redisInstance = &redisService{db: client}
	return redisInstance
}

func (r *redisService) DB() *redis.Client {
	return r.db
}

func (r *redisService) Close() error {
	return r.db.Close()
}
