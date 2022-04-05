package media

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"
	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/handler/httperror"
	"yatter-backend-go/app/handler/utils"
)

func mediatype(contentType string) string {
	if strings.Contains(contentType, "image") {
		return "image"
	}
	if strings.Contains(contentType, "video") {
		return "video"
	}
	return "unknown"
}

func (h *handler) Upload(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	description := r.FormValue("description")

	fileSrc, header, err := r.FormFile("file")
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}
	defer fileSrc.Close()

	mediatype := mediatype(header.Header["Content-Type"][0])
	url := utils.CreateURL(header.Filename)

	attachment := &object.Attachment{
		MediaType:   mediatype,
		URL:         url,
		Description: &description,
	}

	attachment.ID, err = h.app.Dao.Attachment().Insert(ctx, *attachment)
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}

	fileDest, err := os.Create(attachment.URL)
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}
	defer fileDest.Close()

	io.Copy(fileDest, fileSrc)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(attachment); err != nil {
		httperror.InternalServerError(w, err)
		return
	}
}
