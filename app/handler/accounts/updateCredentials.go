package accounts

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/handler/auth"
	"yatter-backend-go/app/handler/files"
	"yatter-backend-go/app/handler/httperror"
)

func uploadMedia(r *http.Request, key string) (*string, error) {
	fileSrc, fileHeader, err := r.FormFile(key)
	if err != nil {
		return nil, err
	}
	defer fileSrc.Close()
	url := files.CreateURL(fileHeader.Filename)
	fileDest, err := os.Create(url)
	if err != nil {
		return nil, err
	}
	defer fileDest.Close()
	io.Copy(fileDest, fileSrc)
	return &url, err
}

func updateObject(r *http.Request, a *object.Account) error {
	displayName := r.FormValue("display_name")
	a.DisplayName = &displayName
	note := r.FormValue("note")
	a.Note = &note

	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		return err
	}

	for k := range r.MultipartForm.File {
		if k == "avatar" {
			a.Avatar, err = uploadMedia(r, k)
			if err != nil {
				return err
			}
		}
		if k == "header" {
			a.Header, err = uploadMedia(r, k)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (h *handler) UpdateCredentials(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// ログインユーザーを取得
	login := auth.AccountOf(r)
	if login == nil {
		httperror.InternalServerError(w, fmt.Errorf("lost account"))
	}

	// 入力内容を取得
	if err := updateObject(r, login); err != nil {
		httperror.InternalServerError(w, err)
		return
	}

	err := h.app.Dao.Account().Update(ctx, *login)
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(login); err != nil {
		httperror.InternalServerError(w, err)
		return
	}
}
