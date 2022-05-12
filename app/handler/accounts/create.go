package accounts

import (
	"encoding/json"
	"fmt"
	"net/http"

	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/handler/httperror"
)

// Request body for "POST /v1/accounts"
type AddRequest struct {
	Username string
	Password string
}

// Handle request for "POST /v1/accounts"
func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// リクエストの内容を取得
	var req AddRequest
	d := json.NewDecoder(r.Body)
	if err := d.Decode(&req); err != nil {
		httperror.BadRequest(w, err)
		return
	}

	if len(req.Username) < 1 {
		httperror.BadRequest(w, fmt.Errorf("empty username"))
		return
	}

	// 同じユーザー名がいるかチェック
	a, err := h.app.Dao.Account().FindByUsername(ctx, req.Username)
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	} else if a != nil {
		httperror.Error(w, http.StatusConflict)
		return
	}

	account := new(object.Account)
	account.Username = req.Username
	if err = account.SetPassword(req.Password); err != nil {
		httperror.InternalServerError(w, err)
		return
	}

	// データベースにアカウント作成
	if _, err = h.app.Dao.Account().Insert(ctx, *account); err != nil {
		httperror.InternalServerError(w, err)
		return
	}

	// データベース上のアカウント情報を取得
	entity, err := h.app.Dao.Account().FindByUsername(ctx, account.Username)
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}

	// アカウント情報をjsonにエンコード
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(entity); err != nil {
		httperror.InternalServerError(w, err)
		return
	}
}
