package status

import (
	"encoding/json"
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
	var req Status
	d := json.NewDecoder(r.Body)
	if err := d.Decode(&req); err != nil {
		httperror.BadRequest(w, err)
		return
	}

	status := new(object.Status)
	status.Content = req.Status
	status.Account = auth.AccountOf(r)

	ctx := r.Context()
	repo := h.app.Dao.Status()

	err := repo.Post(ctx, status)
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
