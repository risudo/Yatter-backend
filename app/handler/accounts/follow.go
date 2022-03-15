package accounts

import (
	"net/http"
	"strings"
	"yatter-backend-go/app/handler/httperror"

	"github.com/go-chi/chi"
)

// Handle request for "POST /v1/accounts/{username}/follow"
func (h *handler) Follow(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	a := r.Header.Get("Authentication")
	pair := strings.SplitN(a, " ", 2)
	if len(pair) < 2 {
		httperror.Error(w, http.StatusUnauthorized)
		return
	}

	//TODO: リファクタリング
	username := chi.URLParam(r, "username")
	arepo := h.app.Dao.Account()
	follower, err := arepo.FindByUsername(ctx, username)
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}
	if follower == nil {
		httperror.Error(w, 404)
		return
	}
	followee, err := arepo.FindByUsername(ctx, pair[1])
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}

	frepo := h.app.Dao.Relation()
	err = frepo.Follow(ctx, followee.ID, follower.ID)
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}
}
