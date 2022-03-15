package accounts

import (
	"net/http"
	"strings"
	"yatter-backend-go/app/handler/httperror"

	"github.com/go-chi/chi"
)

// Handle request for "POST /v1/accounts/{username}/follow"
func (h *handler) Follow(w http.ResponseWriter, r *http.Request) {
	//TODO: リファクタリング
	ctx := r.Context()
	arepo := h.app.Dao.Account()

	a := r.Header.Get("Authentication")
	pair := strings.SplitN(a, " ", 2)
	if len(pair) < 2 {
		httperror.Error(w, http.StatusUnauthorized)
		return
	}

	following, err := arepo.FindByUsername(ctx, pair[1])
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}

	username := chi.URLParam(r, "username")
	follower, err := arepo.FindByUsername(ctx, username)
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}
	if follower == nil {
		httperror.Error(w, 404)
		return
	}
	frepo := h.app.Dao.Relation()
	err = frepo.Follow(ctx, following.ID, follower.ID)
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}
	//TODO: レスポンスに書き込む
}
