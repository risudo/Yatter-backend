package accounts

import (
	"net/http"
	"yatter-backend-go/app/handler/httperror"
	"yatter-backend-go/app/handler/request"

	"github.com/go-chi/chi"
)

func (h *handler) Follow(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := request.IDOf(r)
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}

	username := chi.URLParam(r, "username")
	account := h.app.Dao.Account()
	entity, err := account.FindByUsername(ctx, username)
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}
	if account == nil {
		httperror.Error(w, 404)
		return
	}

	follows := h.app.Dao.Follows()
	follows.Follow(ctx, id, entity)
}
