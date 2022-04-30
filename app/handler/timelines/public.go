package timelines

import (
	"encoding/json"
	"net/http"
	"yatter-backend-go/app/handler/httperror"
	"yatter-backend-go/app/handler/parameters"
)

// Handler request for "GET /v1/timelines/public"
func (h *handler) Public(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	p, err := parameters.ParseAll(r)
	if err != nil {
		httperror.BadRequest(w, err)
		return
	}

	timeline, err := h.app.Dao.Status().PublicTimeline(ctx, *p)
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}

	for i := range timeline {
		timeline[i].MediaAttachments, err = h.app.Dao.Attachment().FindByStatusID(ctx, timeline[i].ID)
		if err != nil {
		httperror.InternalServerError(w, err)
		return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(timeline); err != nil {
		httperror.InternalServerError(w, err)
		return
	}
}
