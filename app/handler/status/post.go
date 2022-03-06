package status

import (
	"encoding/json"
	"log"
	"net/http"
	"yatter-backend-go/app/handler/httperror"
)

type Status struct {
	Status string
}

func (h *handler) Post(w http.ResponseWriter, r *http.Request) {
	var status Status
	d := json.NewDecoder(r.Body)
	if err := d.Decode(&status); err != nil {
		httperror.BadRequest(w, err)
		return
	}
	log.Println(status)

	log.Println(d)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(status); err != nil {
		httperror.InternalServerError(w, err)
		return
	}
}
