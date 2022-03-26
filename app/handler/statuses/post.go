package statuses

import (
	"encoding/json"
	"errors"
	"net/http"
	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/handler/auth"
	"yatter-backend-go/app/handler/httperror"
)

// Requuest body for `POST /v1/statuses`
type Status struct {
	Status string
}

// Handle request for `POST /v1/statuses`
func (h *handler) Post(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req Status
	d := json.NewDecoder(r.Body)
	if err := d.Decode(&req); err != nil {
		httperror.BadRequest(w, err)
		return
	}

	status := &object.Status{
		Content: req.Status,
		Account: auth.AccountOf(r),
	}
	if status.Account == nil {
		httperror.InternalServerError(w, errors.New("lost account"))
		return
	}

	status, err := h.app.Dao.Status().Post(ctx, status)
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(status); err != nil {
		httperror.InternalServerError(w, err)
		return
	}
}
