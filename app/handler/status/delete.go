package status

import (
	"net/http"
	"yatter-backend-go/app/handler/httperror"
	"yatter-backend-go/app/handler/request"
)

// Handler request for "DELETE /v1/statuses/id"
func (h *handler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := request.IDOf(r)
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}

	if err = h.app.Dao.Status().Delete(ctx, id); err != nil {
		httperror.InternalServerError(w, err)
		return
	}
}
