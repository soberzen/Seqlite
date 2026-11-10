package repository

import (
	"context"
	"errors"
	"fmt"
	"qingmo-server/internal/errs"
	"time"

	"github.com/redis/go-redis/v9"
)

type RefreshSessionRedisRepository struct {
	redis *redis.Client
}

type RefreshSession struct {
	UserID    uint      `json:"user_id"`
	JTI       string    `json:"jti"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

const RefreshSessionKey = "refresh_session"

func NewRefreshSessionRedisRepository(
	redis *redis.Client,
) *RefreshSessionRedisRepository {
	return &RefreshSessionRedisRepository{redis: redis}
}

func (r *RefreshSessionRedisRepository) Create(
	ctx context.Context,
	sessionID string,
	refreshSession *RefreshSession,
) error {
	ttl, err := refreshSessionTTL(refreshSession.ExpiresAt)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("%s:%s", RefreshSessionKey, sessionID)

	pipe := r.redis.TxPipeline()
	pipe.HSet(ctx, key,
		"user_id", refreshSession.UserID,
		"jti", refreshSession.JTI,
		"issued_at", refreshSession.IssuedAt.Unix(),
		"expires_at", refreshSession.ExpiresAt.Unix(),
	)
	pipe.PExpire(ctx, key, ttl)
	_, err = pipe.Exec(ctx)
	if errors.Is(err, context.Canceled) ||
		errors.Is(err, context.DeadlineExceeded) {
		return err
	}

	return errs.Wrap(errs.KindInternal, "fail create refresh session", err)
}

func refreshSessionTTL(expiresAt time.Time) (time.Duration, error) {
	ttl := time.Until(expiresAt)
	if ttl <= 0 {
		return 0, errs.New(
			errs.KindInternal,
			"refresh session expiry is in the past",
		)
	}

	if remainder := ttl % time.Millisecond; remainder != 0 {
		ttl += time.Millisecond - remainder
	}
	return ttl, nil
}
