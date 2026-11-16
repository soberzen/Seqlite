package realtime

import (
	"context"
	"net/http"
	"strconv"

	"github.com/Deln0r/ygo/persist"
	"github.com/Deln0r/ygo/server"
)

type documentEditKey struct{}
type actorIDKey struct{}

func WithDocumentAccess(
	r *http.Request,
	documentID int64,
	actorID uint,
	canEdit bool,
) *http.Request {
	ctx := context.WithValue(r.Context(), documentEditKey{}, canEdit)
	ctx = context.WithValue(ctx, actorIDKey{}, actorID)
	req := r.WithContext(ctx)
	req.SetPathValue("documentID", strconv.FormatInt(documentID, 10))
	return req
}

func GetActorID(ctx context.Context) (uint, bool) {
	id, ok := ctx.Value(actorIDKey{}).(uint)
	return id, ok
}

func NewCollabHandler(store persist.Store, origins []string) http.Handler {
	s := server.New(server.Options{
		Store:          store,
		OriginPatterns: origins,
		FlushEvery:     200,
		DocNameFn: func(r *http.Request) string {
			return r.PathValue("documentID")
		},
		ReadOnly: func(_ string, r *http.Request) bool {
			canEdit, _ := r.Context().Value(documentEditKey{}).(bool)
			return !canEdit

		},
	})

	return s.Handler()
}
