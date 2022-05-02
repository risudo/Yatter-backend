package timelines

import (
	"encoding/json"
	"fmt"
	"net/http"
	"yatter-backend-go/app/handler/auth"
	"yatter-backend-go/app/handler/httperror"
	"yatter-backend-go/app/handler/parameters"
)

func (h *handler) Home(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	login := auth.AccountOf(r)
	if login == nil {
		httperror.InternalServerError(w, fmt.Errorf("lost account"))
		return
	}

	p, err := parameters.ParseAll(r)
	if err != nil {
		httperror.BadRequest(w, err)
		return
	}

	timeline, err := h.app.Dao.Status().HomeTimeline(ctx, login.ID, *p)
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
