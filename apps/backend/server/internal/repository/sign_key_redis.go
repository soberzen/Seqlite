package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"qingmo-server/internal/security"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/singleflight"
)

type SignKeyRedisRepository struct {
	redis *redis.Client
	sf    singleflight.Group
}

func NewSignKeyRedisRepository(client *redis.Client) *SignKeyRedisRepository {
	return &SignKeyRedisRepository{redis: client}
}

const (
	signingPeriod    = 30 * 24 * time.Hour
	operationTimeout = 5 * time.Second
	maxKeyAttempts   = 5
)

// keys[1] active kid
// keys[2] key path
// argv[1] expected active kid
// argv[2] retainUntil
var checkAndExtendSigningKeyScript = redis.NewScript(`
-- 检查当前活动 kid 是否存在，不存在说明密钥轮换过
local currentActiveKid = redis.call("GET", KEYS[1])
if not currentActiveKid or currentActiveKid == "" then
	return {"retry"}
end

-- 检查当前活动 kid 是否与预期一致，如果不是的话就说明密钥轮换过
if currentActiveKid ~= ARGV[1] then
	return {"retry"}
end

-- 活动密钥存在且与预期一致
local clock = redis.call("TIME")
local now = tonumber(clock[1]) * 1000 + math.floor(tonumber(clock[2]) / 1000)
local retainUntil = tonumber(ARGV[2])

if not retainUntil or retainUntil <= now then
 return {"invalid_retain_until"}
end
local privateKey = redis.call("HGET", KEYS[2], "private_key")
local publicKey = redis.call("HGET", KEYS[2], "public_key")
local rotateAt = tonumber(redis.call("HGET", KEYS[2], "rotate_at"))

-- 检查密钥是否存在且未过期
if not privateKey or privateKey == "" or
   not publicKey or publicKey == "" or
   not rotateAt or rotateAt <= now then
    return {"rotate"}
end

local ttl = redis.call("PTTL", KEYS[2])
if ttl == -2 then
	return {"retry"}
end

-- 检查密钥是否过期, 如果过期则刷新过期时间
-- 就是为了解决当前密钥对并没有到过期时间，但是在jwt签名后，密钥对过期了，延长过期时间到 retainUntil
if ttl >= 0 and now + ttl < retainUntil then
    redis.call("PEXPIREAT", KEYS[2], retainUntil)
end
return {"ok", privateKey, publicKey}
`)

// CAS 轮换：只有活动 kid 未改变时才发布新密钥，不修改旧密钥。
var publishSignKeyScript = redis.NewScript(`
local current = redis.call("GET", KEYS[1]) or ""
if current ~= ARGV[1] then
    return 0
end
if redis.call("EXISTS", KEYS[2]) == 1 then
    return -1
end

local period = tonumber(ARGV[5])
if not period or period < 1 then
    return redis.error_reply("invalid signing period")
end

local clock = redis.call("TIME")
local now = tonumber(clock[1]) * 1000 + math.floor(tonumber(clock[2]) / 1000)
local rotateAt = now + period

redis.call("HSET", KEYS[2],
    "private_key", ARGV[3],
    "public_key", ARGV[4],
    "rotate_at", rotateAt)
redis.call("PEXPIREAT", KEYS[2], rotateAt)
redis.call("SET", KEYS[1], ARGV[2])
return 1
`)

func keyPath(t security.TokenType, kid string) string {
	return fmt.Sprintf("auth:rsa:%s:%s", t, kid)
}

func activePath(t security.TokenType) string {
	return fmt.Sprintf("auth:rsa:active:%s", t)
}

func (r *SignKeyRedisRepository) GetActiveKey(
	ctx context.Context,
	tokenType security.TokenType,
	retainUntil time.Time,
) (string, *security.RSAKeyPair, error) {
	ctx, cancel := context.WithTimeout(ctx, operationTimeout)
	defer cancel()

	if !retainUntil.After(time.Now()) {
		return "", nil, fmt.Errorf(
			"key retention deadline must be in the future",
		)
	}

	for range maxKeyAttempts {
		if err := ctx.Err(); err != nil {
			return "", nil, err
		}
		// 查看当前活动指针。
		kid, err := r.redis.Get(ctx, activePath(tokenType)).Result()
		// 没有活动指针，创建新密钥。
		if errors.Is(err, redis.Nil) {
			if err := r.ensureNewKey(ctx, tokenType, ""); err != nil {
				return "", nil, err
			}
			continue
		}
		if err != nil {
			return "", nil, fmt.Errorf("get active RSA kid: %w", err)
		}
		// 确认当前 kid、检查签发期限，并保证保留到 retainUntil。
		result, err := checkAndExtendSigningKeyScript.Run(
			ctx,
			r.redis,
			[]string{activePath(tokenType), keyPath(tokenType, kid)},
			kid,
			retainUntil.UnixMilli(),
		).Slice()
		if err != nil {
			return "", nil, fmt.Errorf("claim signing key: %w", err)
		}
		if len(result) == 0 {
			return "", nil, fmt.Errorf("empty signing key result")
		}

		status, ok := result[0].(string)
		if !ok {
			return "", nil, fmt.Errorf("invalid signing key status")
		}

		switch status {
		case "ok":
			if err := ctx.Err(); err != nil {
				return "", nil, err
			}
			if len(result) != 3 {
				return "", nil, fmt.Errorf("invalid signing key result length")
			}

			privateKey, privateOK := result[1].(string)
			publicKey, publicOK := result[2].(string)
			if !privateOK || !publicOK || privateKey == "" || publicKey == "" {
				return "", nil, fmt.Errorf("invalid signing key pair")
			}

			return kid, &security.RSAKeyPair{
				PrivatePEM: privateKey,
				PublicPEM:  publicKey,
			}, nil

		case "rotate":
			if err := r.ensureNewKey(ctx, tokenType, kid); err != nil {
				return "", nil, err
			}

		case "retry":
			// 活动指针刚被其他请求更新，下一轮重新读取。
		case "invalid_retain_until":
			return "", nil, fmt.Errorf("key retention deadline has passed")

		default:
			return "", nil, fmt.Errorf("unknown signing key status: %q", status)
		}
	}

	return "", nil, fmt.Errorf(
		"signing key changed too often; retry limit reached",
	)
}

// singleflight 只合并创建工作；每个调用者仍须单独执行领取脚本。
func (r *SignKeyRedisRepository) ensureNewKey(
	ctx context.Context,
	tokenType security.TokenType,
	expectedKid string,
) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	flightKey := fmt.Sprintf("publish:%s:%s", tokenType, expectedKid)

	resultCh := r.sf.DoChan(flightKey, func() (any, error) {
		// 共享工作使用独立超时，避免第一个请求取消后影响其他等待者。
		workCtx, cancel := context.WithTimeout(
			context.WithoutCancel(ctx),
			operationTimeout,
		)
		defer cancel()

		for range maxKeyAttempts {
			if err := workCtx.Err(); err != nil {
				return nil, err
			}

			kid := uuid.New().String()

			pair, err := security.GenerateRSAPair()
			if err != nil {
				return nil, fmt.Errorf("generate RSA key: %w", err)
			}

			if err := workCtx.Err(); err != nil {
				return nil, err
			}

			published, err := publishSignKeyScript.Run(
				workCtx,
				r.redis,
				[]string{activePath(tokenType), keyPath(tokenType, kid)},
				expectedKid,
				kid,
				pair.PrivatePEM,
				pair.PublicPEM,
				signingPeriod.Milliseconds(),
			).Int()
			if err != nil {
				return nil, fmt.Errorf("publish RSA key: %w", err)
			}

			switch published {
			case 0, 1:
				// 自己发布成功，或其他请求已经完成轮换。
				return nil, nil

			case -1:
				// 新 kid 碰撞，重新生成。
				continue

			default:
				return nil, fmt.Errorf(
					"invalid key publish result: %d",
					published,
				)
			}
		}

		return nil, fmt.Errorf("RSA kid collision retry limit reached")
	})

	// 同时等待共享结果和当前请求取消。
	select {
	case <-ctx.Done():
		return ctx.Err()

	case result := <-resultCh:
		if err := ctx.Err(); err != nil {
			return err
		}
		return result.Err
	}
}

func (r *SignKeyRedisRepository) Rotate(
	ctx context.Context,
	tokenType security.TokenType,
	oldKid string,
	retainUntil time.Time,
) (string, *security.RSAKeyPair, error) {
	ctx, cancel := context.WithTimeout(ctx, operationTimeout)
	defer cancel()

	if !retainUntil.After(time.Now()) {
		return "", nil, fmt.Errorf(
			"key retention deadline must be in the future",
		)
	}

	if err := r.ensureNewKey(ctx, tokenType, oldKid); err != nil {
		return "", nil, err
	}

	// 轮换后仍要用本次请求的 retainUntil 领取密钥。
	return r.GetActiveKey(ctx, tokenType, retainUntil)
}

func (r *SignKeyRedisRepository) GetKeyPair(
	ctx context.Context,
	tokenType security.TokenType,
	kid string,
) (*security.RSAKeyPair, error) {
	ctx, cancel := context.WithTimeout(ctx, operationTimeout)
	defer cancel()

	// 验签按 JWT 中的 kid 查找，允许读取已经停止签发的旧密钥。
	values, err := r.redis.HMGet(
		ctx,
		keyPath(tokenType, kid),
		"private_key",
		"public_key",
	).Result()
	if err != nil {
		return nil, fmt.Errorf("load RSA key: %w", err)
	}

	if len(values) != 2 {
		return nil, fmt.Errorf("invalid stored RSA key result")
	}

	privateKey, privateOK := values[0].(string)
	publicKey, publicOK := values[1].(string)
	if !privateOK || !publicOK || privateKey == "" || publicKey == "" {
		return nil, redis.Nil
	}

	return &security.RSAKeyPair{
		PrivatePEM: privateKey,
		PublicPEM:  publicKey,
	}, nil
}
