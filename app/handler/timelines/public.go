package timelines

import (
	"encoding/json"
	"net/http"
	"yatter-backend-go/app/handler/httperror"
)

func (h *handler) Public(w http.ResponseWriter, r *http.Request) {
	//TODO: statusのスライスを作ればよさそう
	ctx := r.Context()
	repo := h.app.Dao.Status()
	timeline, err := repo.PublicTimeline(ctx)
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
