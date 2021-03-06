package accounts

import (
	"encoding/json"
	"net/http"
	"yatter-backend-go/app/handler/httperror"
	"yatter-backend-go/app/handler/parameters"

	"github.com/go-chi/chi"
)

// Handle request for "GET /v1/accounts/username/followers"
func (h *handler) Followers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	account, err := h.app.Dao.Account().FindByUsername(ctx, chi.URLParam(r, "username"))
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}
	if account == nil {
		httperror.Error(w, http.StatusNotFound)
		return
	}

	p, err := parameters.ParseAll(r)
	if err != nil {
		httperror.BadRequest(w, err)
		return
	}

	accounts, err := h.app.Dao.Relation().Followers(ctx, account.ID, *p)
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(accounts); err != nil {
		httperror.InternalServerError(w, err)
		return
	}
}
