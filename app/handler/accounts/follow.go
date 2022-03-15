package accounts

import (
	"net/http"
	"yatter-backend-go/app/handler/auth"
	"yatter-backend-go/app/handler/httperror"

	"github.com/go-chi/chi"
)

// Handle request for "POST /v1/accounts/{username}/follow"
func (h *handler) Follow(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	following := auth.AccountOf(r)
	if following != nil {
		httperror.InternalServerError(w, nil)//
		return
	}

	followerName := chi.URLParam(r, "username")
	follower, err := h.app.Dao.Account().FindByUsername(ctx, followerName)
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}
	if follower == nil {
		httperror.Error(w, 404)
		return
	}
	err = h.app.Dao.Relation().Follow(ctx, following.ID, follower.ID)
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}
	//TODO: レスポンスに書き込む
}
