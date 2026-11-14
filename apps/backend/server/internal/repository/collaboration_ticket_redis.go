package repository

import (
	"context"
	"qingmo-server/internal/errs"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type CollaborationTicketRedisRepository struct {
	redis *redis.Client
}

func NewCollaborationTicketRedisRepository(
	redis *redis.Client,
) *CollaborationTicketRedisRepository {
	return &CollaborationTicketRedisRepository{
		redis: redis,
	}
}

const collaborationTicketTTL = 60 * time.Second

func pathKey(ticket string) string {
	return "collaboration_ticket:" + ticket
}

func (h *CollaborationTicketRedisRepository) Save(
	ctx context.Context,
	publicID uuid.UUID,
	ticket string,
	userID uint,
) error {

	key := pathKey(ticket)

	pipe := h.redis.TxPipeline()

	pipe.HSet(ctx, key, "publicID", publicID.String(), "userID", userID)
	pipe.Expire(ctx, key, collaborationTicketTTL)
	_, err := pipe.Exec(ctx)

	if err != nil {
		return errs.Wrap(errs.KindInternal, "fail create ticket", err)
	}
	return nil
}

func (h *CollaborationTicketRedisRepository) Consume(
	ctx context.Context,
	publicID uuid.UUID,
	ticket string,
) (uint, error) {
	key := pathKey(ticket)

	script := redis.NewScript(`
		local data = redis.call("HGETALL", KEYS[1])
		if #data == 0 then
			return nil
		end
		redis.call("DEL", KEYS[1])
		return data
	`)
	vals, err := script.Run(ctx, h.redis, []string{key}, publicID.String()).
		Result()
	if err != nil {
		return 0, errs.Wrap(errs.KindInternal, "failed consume ticket", err)
	}
	if vals == nil {
		return 0, errs.New(errs.KindInvalid, "ticket not found")
	}
	arr, ok := vals.([]interface{})
	if !ok || len(arr)%2 != 0 {
		return 0, errs.New(errs.KindInternal, "corrupt ticket payload")
	}

	var storedPublicID string
	var userIDStr string
	for i := 0; i+1 < len(arr); i += 2 {
		field, _ := arr[i].(string)
		value, _ := arr[i+1].(string)
		switch field {
		case "publicID":
			storedPublicID = value
		case "userID":
			userIDStr = value
		}
	}
	if storedPublicID != publicID.String() {
		return 0, errs.New(errs.KindForbidden, "ticket does not match document")
	}
	uid, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		return 0, errs.Wrap(errs.KindInternal, "corrupt userID in ticket", err)
	}
	if uid == 0 {
		return 0, errs.New(errs.KindInternal, "corrupt userID = 0 in ticket")
	}
	return uint(uid), nil
}
