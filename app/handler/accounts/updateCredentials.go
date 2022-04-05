package accounts

import (
	"fmt"
	"net/http"
	"yatter-backend-go/app/handler/auth"
	"yatter-backend-go/app/handler/httperror"
	"yatter-backend-go/app/handler/utils"
)

/*
curl -X 'POST' \
  'http://localhost:8080/v1/accounts/update_credentials' \
  -H 'accept: application/json' \
  -H 'Content-Type: multipart/form-data' \
  -F 'display_name=' \
  -F 'note=' \
  -F 'avatar=@Screen Shot 2022-04-04 at 9.06.05.png;type=image/png' \
  -F 'header=@Screen Shot 2022-02-15 at 17.27.11.png;type=image/png'
*/

func (h *handler) UpdateCredentials(w http.ResponseWriter, r *http.Request) {
	// ログインユーザーを取得
	login := auth.AccountOf(r)
	if login == nil {
		httperror.InternalServerError(w, fmt.Errorf("lost account"))
	}

	// 入力内容を取得
	avatarSrc, avatarHeader, err := r.FormFile("avatar")
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}
	defer avatarSrc.Close()
	avatarURL := utils.CreateURL(avatarHeader.Filename)

	headerSrc, headerHeader, err := r.FormFile("header")
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}
	defer headerSrc.Close()
	headerURL := utils.CreateURL(headerHeader.Filename)

	// 更新
	login.Avatar = &avatarURL
	login.Header = &headerURL


}
