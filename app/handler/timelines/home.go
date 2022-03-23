package timelines

import (
	"encoding/json"
	"errors"
	"net/http"
	"yatter-backend-go/app/handler/auth"
	"yatter-backend-go/app/handler/httperror"
)

func (h *handler) Home(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	login := auth.AccountOf(r)
	if login == nil {
		httperror.InternalServerError(w, errors.New("lost account"))
		return
	}

	timeline, err := h.app.Dao.Status().HomeTimeline(ctx, login.ID)
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
