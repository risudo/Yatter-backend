package accounts

import (
	"context"
	"encoding/json"
	"net/http"

	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/domain/repository"
	"yatter-backend-go/app/handler/httperror"
)

// Request body for "POST /v1/accounts"
type AddRequest struct {
	Username     string
	Password     string
	Display_Name string
	Note         string
	Avatar       string
	Header       string
}

// custom error type
type (
	errBadrequest struct {
		message string
	}

	errConflict struct{}
)

func (e *errBadrequest) Error() string {
	return e.message
}

func (e *errConflict) Error() string {
	return ""
}

/*
Example Request body

{
   "username":"john",
   "password":"P@ssw0rd",
   "avatar":"attachments/2022-05-13T02:14:24.4652768Z.png",
   "note":"note",
   "header":"attachments/2022-05-13T02:14:24.4652768Z.png",
   "display_name":"display "
}
*/

func parseRequest(ctx context.Context, r *http.Request, repo repository.Account) (*object.Account, error) {
	var req AddRequest
	d := json.NewDecoder(r.Body)
	if err := d.Decode(&req); err != nil {
		return nil, &errBadrequest{message: err.Error()}
	}

	if len(req.Username) < 1 {
		return nil, &errBadrequest{message: "empty username"}
	}

	// 同じユーザー名がいるかチェック
	a, err := repo.FindByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	} else if a != nil {
		return nil, &errConflict{}
	}

	account := &object.Account{
		Username: req.Username,
	}

	if err = account.SetPassword(req.Password); err != nil {
		return nil, err
	}
	return account, nil
}

// Handle request for "POST /v1/accounts"
func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	account, err := parseRequest(ctx, r, h.app.Dao.Account())
	if err != nil {
		switch err.(type) {
		case *errBadrequest:
			httperror.BadRequest(w, err)
			return
		case *errConflict:
			httperror.Error(w, http.StatusConflict)
			return
		default:
			httperror.InternalServerError(w, err)
			return
		}
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
