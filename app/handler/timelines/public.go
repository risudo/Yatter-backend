package timelines

import (
	"encoding/json"
	"net/http"
	"yatter-backend-go/app/handler/httperror"
)

// Handler request for "GET /v1/timelines/public"
func (h *handler) Public(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	timeline, err := h.app.Dao.Status().PublicTimeline(ctx)
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(timeline); err != nil {
		httperror.InternalServerError(w, err)
		return
	}
}
