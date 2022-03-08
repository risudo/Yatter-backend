package status

import (
	"encoding/json"
	"net/http"
	"yatter-backend-go/app/handler/httperror"
	"yatter-backend-go/app/handler/request"
)

func (h *handler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	repo := h.app.Dao.Status()
	id, err := request.IDOf(r)
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}

	status, err := repo.FindById(ctx, id)

	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}
	if status == nil {
		httperror.Error(w, 404)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(status); err != nil {
		httperror.InternalServerError(w, err)
		return
	}
}
