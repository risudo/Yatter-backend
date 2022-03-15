package accounts

import (
	"encoding/json"
	"net/http"
	"yatter-backend-go/app/handler/httperror"

	"github.com/go-chi/chi"
)

// Handle request for "GET /v1/accounts/username/followers"
func (h *handler) Followers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	username := chi.URLParam(r, "username")
	arepo := h.app.Dao.Account()
	account, err := arepo.FindByUsername(ctx, username)
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}
	if account == nil {
		httperror.Error(w, 404)
		return
	}

	frepo := h.app.Dao.Relation()
	accounts, err := frepo.Followers(ctx, account.ID)
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}
	if accounts == nil {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(accounts); err != nil {
		httperror.InternalServerError(w, err)
		return
	}
}
