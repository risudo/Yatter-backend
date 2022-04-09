package statuses

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/handler/auth"
	"yatter-backend-go/app/handler/httperror"
)


// Handle request for `POST /v1/statuses`
func (h *handler) Post(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}
	var jsonBody map[string]string
	err = json.Unmarshal(body, &jsonBody)
	if err != nil {
		httperror.BadRequest(w, err)
		return
	}

	status := &object.Status{
		Content: jsonBody["status"],
		Account: auth.AccountOf(r),
	}
	if status.Account == nil {
		httperror.InternalServerError(w, fmt.Errorf("lost account"))
		return
	}

	id, err := h.app.Dao.Status().Insert(ctx, *status)
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}

	entity, err := h.app.Dao.Status().FindByID(ctx, id)
	if err != nil || entity == nil {
		httperror.InternalServerError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(entity); err != nil {
		httperror.InternalServerError(w, err)
		return
	}
}
