package statuses

import (
	"encoding/json"
	"fmt"
	"net/http"
	"yatter-backend-go/app/handler/auth"
	"yatter-backend-go/app/handler/httperror"
	"yatter-backend-go/app/handler/request"
)

// Handler request for "DELETE /v1/statuses/id"
func (h *handler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := request.IDOf(r)
	if err != nil {
		httperror.BadRequest(w, err)
		return
	}

	status, err := h.app.Dao.Status().FindByID(ctx, id)
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}
	if status == nil {
		httperror.Error(w, http.StatusNotFound)
		return
	}

	// statusの投稿者とログインユーザーの一致を確認
	login := auth.AccountOf(r)
	if login == nil {
		httperror.InternalServerError(w, fmt.Errorf("lost account"))
		return
	}
	if status.Account.ID != login.ID {
		httperror.BadRequest(w, fmt.Errorf("status does not belong to the User"))
		return
	}

	if err = h.app.Dao.Status().Delete(ctx, id); err != nil {
		httperror.InternalServerError(w, err)
		return
	}

	if err := json.NewEncoder(w).Encode(&struct{}{}); err != nil {
		httperror.InternalServerError(w, err)
		return
	}
}
