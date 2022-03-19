package accounts

import (
	"encoding/json"
	"log"
	"net/http"
	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/handler/auth"
	"yatter-backend-go/app/handler/httperror"
)

// Handler request for "GET /v1/accounts/Relationships"
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
	} else if account == nil {
		httperror.Error(w, http.StatusNotFound)
		return
	}

	// TODO: Followと被ってるからまとめた方がいい？
	relationRepo := h.app.Dao.Relation()
	relation := new(object.RelationWith)
	relation.Following, err = relationRepo.IsFollowing(ctx, requester.ID, account.ID)
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}
	relation.FollowedBy, err = relationRepo.IsFollowing(ctx, account.ID, requester.ID)
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
