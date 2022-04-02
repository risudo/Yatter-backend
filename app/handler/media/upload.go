package media

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"yatter-backend-go/app/handler/httperror"
)

func (h *handler) Upload(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fileSrc, _, err := r.FormFile("file")
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}
	defer fileSrc.Close()
	fileDest, err := os.Create("tmp.png")
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}
	defer fileDest.Close()

	io.Copy(fileDest, fileSrc)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(&struct{}{}); err != nil {
		httperror.InternalServerError(w, err)
		return
	}
}
