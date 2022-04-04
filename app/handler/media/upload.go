package media

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"
	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/handler/httperror"
)

// curl -F 'file=@Screen Shot 2022-04-01 at 17.42.58.png;type=image/png'

var nowid int64 = 1

func (h *handler) Upload(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	description := r.FormValue("description")

	fileSrc, header, err := r.FormFile("file")
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}
	defer fileSrc.Close()

	mediatype := header.Header["Content-Type"][0]
	var ext string
	switch mediatype {
	case "image/png":
		ext = ".png"
	}
	url := "attachments/" + strconv.FormatInt(nowid, 10) + ext
	nowid++

	attachment := &object.Attachment{
		MediaType:   header.Header["Content-Type"][0],
		URL: url,
		Description: &description,
	}

	attachment.ID, err = h.app.Dao.Attachment().Insert(ctx, *attachment) // IDをリターンする
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

/*
'Content-Type: multipart/form-data'

note: IDとファイルパスの作り方案
1. いったんdbにインサートしてそのIDを拾ってくる
2. 静的変数としてIDを保存してインクリメントしていく

idとファイル名を一致させる？
リターンされるIDがDBを紐づいていることは保証されていた方が良さそう
*/
