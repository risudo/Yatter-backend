package accounts

import (
	"encoding/json"
	"log"
	"net/http"
	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/handler/auth"
	"yatter-backend-go/app/handler/httperror"
)

func (h *handler) Relationships(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	requester := auth.AccountOf(r)
	if requester == nil {
		log.Println("auth fail") //TODO:ちゃんとエラーを定義する
		return
	}

	username := r.FormValue("username")

	account, err := h.app.Dao.Account().FindByUsername(ctx, username)
	if err != nil {
		httperror.InternalServerError(w, err)
		return
		// アカウントが見つからなかった場合
	} else if account == nil {
		return
	}

	// TODO: Followと被ってるからまとめたい
	repo := h.app.Dao.Relation()
	relation := new(object.RelationWith)
	relation.Following, err = repo.IsFollowing(ctx, requester.ID, account.ID)
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}
	relation.FollowedBy, err = repo.IsFollowing(ctx, account.ID, requester.ID)
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(relation); err != nil {
		httperror.InternalServerError(w, err)
		return
	}
}
