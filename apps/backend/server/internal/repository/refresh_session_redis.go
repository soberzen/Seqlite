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

var rotateScript = redis.NewScript(`
local current = redis.call('HGET', KEYS[1], 'jti')
if not current then
	return 0
end
if current ~= ARGV[1] then
	return -1
end
redis.call('DEL', KEYS[1])
redis.call("HSET", KEYS[2],
	"user_id", ARGV[2],
	"jti", ARGV[3],
	"issued_at", ARGV[4],
	"expires_at", ARGV[5]
)
redis.call('PEXPIRE', KEYS[2], ARGV[6])
return 1
`)

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
	if err != nil {
		return errs.Wrap(errs.KindInternal, "fail create refresh session", err)
	}
	return nil
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

func (r *RefreshSessionRedisRepository) Delete(
	ctx context.Context,
	sessionID string,
) error {
	key := fmt.Sprintf("%s:%s", RefreshSessionKey, sessionID)

	err := r.redis.Del(ctx, key).Err()

	if errors.Is(err, context.Canceled) ||
		errors.Is(err, context.DeadlineExceeded) {
		return err
	}

	if err != nil {
		return errs.Wrap(errs.KindInternal, "fail delete refresh session", err)
	}

	return nil
}

func (r *RefreshSessionRedisRepository) Rotate(
	ctx context.Context,
	oldSessionID,
	expecteOldJTI,
	newSessionID string,
	newRefreshSession *RefreshSession,
) error {
	ttl, err := refreshSessionTTL(newRefreshSession.ExpiresAt)
	if err != nil {
		return err
	}
	oldKey := fmt.Sprintf("%s:%s", RefreshSessionKey, oldSessionID)
	newKey := fmt.Sprintf("%s:%s", RefreshSessionKey, newSessionID)

	result, err := rotateScript.Run(
		ctx,
		r.redis,
		[]string{oldKey, newKey},
		expecteOldJTI,
		newRefreshSession.UserID,
		newRefreshSession.JTI,
		newRefreshSession.IssuedAt.Unix(),
		newRefreshSession.ExpiresAt.Unix(),
		ttl.Milliseconds(),
	).Int64()
	if err != nil {
		if errors.Is(err, context.Canceled) ||
			errors.Is(err, context.DeadlineExceeded) {
			return err
		}
		return errs.Wrap(errs.KindInternal, "fail rotate refresh session", err)
	}
	switch result {
	case 1:
		return nil
	case 0:
		// 旧会话不存在：过期或已登出
		return errs.New(errs.KindNotFound, "refresh session not found")
	case -1:
		// jti 不匹配：已被并发轮换或重放
		return errs.New(errs.KindConflict, "refresh session already rotated")
	default:
		return errs.New(errs.KindInternal, "unexpected rotate result")
	}
}
