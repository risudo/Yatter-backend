package status

import (
	"net/http"
	"strconv"
	"yatter-backend-go/app/handler/httperror"

	"github.com/go-chi/chi"
)

func (h *handler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	repo := h.app.Dao.Status()

	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
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
