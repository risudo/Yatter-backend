package accounts

import (
	"net/http"
	"yatter-backend-go/app/handler/httperror"

	"github.com/go-chi/chi"
)

func (h *handler) Following(w http.ResponseWriter, r *http.Request) {
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
	_, err = frepo.Following(ctx, account.ID)
}
