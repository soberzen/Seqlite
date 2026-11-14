package realtime

import (
	"context"
	"net/http"
	"strconv"

	"github.com/Deln0r/ygo/server"
)

type documentEditKey struct{}

func WithDocumentAccess(
	r *http.Request,
	documentID int64,
	canEdit bool,
) *http.Request {
	ctx := context.WithValue(r.Context(), documentEditKey{}, canEdit)
	req := r.WithContext(ctx)
	req.SetPathValue("documentID", strconv.FormatInt(documentID, 10))
	return req
}

func NewCollabHandler(origins []string) http.Handler {
	s := server.New(server.Options{
		// Store:          store,
		OriginPatterns: origins,
		FlushEvery:     100,
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
