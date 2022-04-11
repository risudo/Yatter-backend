package statuses

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/handler/auth"
	"yatter-backend-go/app/handler/httperror"
)

type AddRequest struct {
	Status    string
	Media_ids []object.AttachmentID
}

// Handle request for `POST /v1/statuses`
func (h *handler) Post(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req AddRequest
	d := json.NewDecoder(r.Body)
	if err := d.Decode(&req); err != nil {
		httperror.BadRequest(w, err)
		return
	}

	log.Println("media_ids:", req.Media_ids)

	status := &object.Status{
		Content: req.Status,
		Account: auth.AccountOf(r),
	}
	if status.Account == nil {
		httperror.InternalServerError(w, fmt.Errorf("lost account"))
		return
	}

	id, err := h.app.Dao.Status().Insert(ctx, *status, req.Media_ids)
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}

	entity, err := h.app.Dao.Status().FindByID(ctx, id)
	if err != nil || entity == nil {
		httperror.InternalServerError(w, err)
		return
	}
	entity.MediaAttachiments, err = h.app.Dao.Attachment().FindByStatusID(ctx, id)
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(entity); err != nil {
		httperror.InternalServerError(w, err)
		return
	}
}
