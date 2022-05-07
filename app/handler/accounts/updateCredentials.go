package accounts

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/handler/auth"
	"yatter-backend-go/app/handler/files"
	"yatter-backend-go/app/handler/httperror"
)

func uploadMedia(r *http.Request, key string) (*string, error) {
	src, fileHeader, err := r.FormFile(key)
	if err != nil {
		return nil, err
	}
	defer func() {
		err := src.Close()
		if err != nil {
			log.Println("Close:", err)
		}
	}()

	err = files.MightCreateAttachmentDir()
	if err != nil {
		return nil, err
	}
	url := files.CreateURL(fileHeader.Filename)
	dest, err := os.Create(url)
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}
	defer func() {
		err := dest.Close()
		if err != nil {
			log.Println("Close:", err)
		}
	}()
	_, err = io.Copy(dest, src)
	if err != nil {
		return nil, err
	}
	return &url, err
}

func updateObject(r *http.Request, a *object.Account) error {
	displayName := r.FormValue("display_name")
	a.DisplayName = &displayName
	note := r.FormValue("note")
	a.Note = &note

	const maxMemory = 32 << 20
	err := r.ParseMultipartForm(maxMemory)
	if err != nil {
		return fmt.Errorf("ParseMultipartForm: %w", err)
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

// Handle request for "POST /v1/update_credentials"
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
