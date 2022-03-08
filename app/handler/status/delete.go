package status

import (
	"net/http"
	"yatter-backend-go/app/handler/httperror"
	"yatter-backend-go/app/handler/request"
)

func (h *handler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	repo := h.app.Dao.Status()

	id, err := request.IDOf(r)
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}
	err = repo.Delete(ctx, id)
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}
}
