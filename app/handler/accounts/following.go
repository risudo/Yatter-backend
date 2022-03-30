package accounts

import (
	"encoding/json"
	"net/http"
	"yatter-backend-go/app/handler/httperror"
	"yatter-backend-go/app/handler/parameters"

	"github.com/go-chi/chi"
)

// Handle request for "GET /v1/accounts/username/following"
func (h *handler) Following(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	username := chi.URLParam(r, "username")
	account, err := h.app.Dao.Account().FindByUsername(ctx, username)
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}
	if account == nil {
		httperror.Error(w, http.StatusNotFound)
		return
	}

	p, err := parameters.Parse(r)
	if err != nil {
		httperror.Error(w, http.StatusBadRequest)
		return
	}

	accounts, err := h.app.Dao.Relation().Following(ctx, account.ID, *p)
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
