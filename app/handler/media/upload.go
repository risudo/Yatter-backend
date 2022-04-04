package media

import (
	"encoding/json"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/handler/httperror"
)

//TODO: シード値ちゃんとしてる？

func randomString(n int) string {
	var letter = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letter[rand.Intn(len(letter))]
	}
	return string(b)
}

func createURL(filename string) string {
	return "attachments/" + time.Now().Format(time.RFC3339Nano) + randomString(5) + filepath.Ext(filename)
}

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
	// mediatype := header.Header["Content-Type"][0]
	url := createURL(header.Filename)

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
