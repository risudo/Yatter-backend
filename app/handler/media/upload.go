package media

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/handler/files"
	"yatter-backend-go/app/handler/httperror"
)

func mediatype(contentType string) string {
	if strings.Contains(contentType, "image") {
		return "image"
	} else if strings.Contains(contentType, "video") {
		return "video"
	} else if strings.Contains(contentType, "gifv") {
		return "gifv"
	}
	return "unknown"
}

// Handle request for "POST /v1/media"
func (h *handler) Upload(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	description := r.FormValue("description")

	src, header, err := r.FormFile("file")
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}
	defer func() {
		if err := src.Close(); err != nil {
			log.Println(err)
		}
	}()

	attachment := &object.Attachment{
		MediaType:   mediatype(header.Header["Content-Type"][0]),
		URL:         files.CreateURL(header.Filename),
		Description: &description,
	}
	if *attachment.Description == "" {
		attachment.Description = nil
	}

	err = files.MightCreateAttachmentDir()
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}

	dest, err := os.Create(attachment.URL)
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}
	defer func() {
		if err := dest.Close(); err != nil {
			log.Println(err)
		}
	}()

	_, err = io.Copy(dest, src)
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}

	attachment.ID, err = h.app.Dao.Attachment().Insert(ctx, *attachment)
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(attachment); err != nil {
		httperror.InternalServerError(w, err)
		return
	}
}
