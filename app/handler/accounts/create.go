package accounts

import (
	"encoding/json"
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

	var req AddRequest
	d := json.NewDecoder(r.Body)
	if err := d.Decode(&req); err != nil {
		httperror.BadRequest(w, err)
		return
	}

	repo := h.app.Dao.Account()

	// 既にユーザーが存在していたら何もしない
	a, err := repo.FindByUsername(ctx, req.Username)
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	} else if a != nil {
		return
	}

	account := new(object.Account)
	account.Username = req.Username
	if err := account.SetPassword(req.Password); err != nil {
		httperror.InternalServerError(w, err)
		return
	}

	//データベースにアカウント作成
	err = repo.Create(ctx, account)
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(account); err != nil {
		httperror.InternalServerError(w, err)
		return
	}
}
