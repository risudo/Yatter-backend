package accounts

import (
	"encoding/json"
	"net/http"
	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/handler/auth"
	"yatter-backend-go/app/handler/httperror"

	"github.com/go-chi/chi"
)

//TODO:フォローしてる人、フォローされる人の変数名もっとわかりやすくしたい

// Handler request for "POST /v1/accounts/usernmae/unfollow"
func (h *handler) Unfollow(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	requester := auth.AccountOf(r)
	if requester == nil {
		httperror.InternalServerError(w, nil) //TODO: ちゃんとエラーを定義する
		return
	}

	followerName := chi.URLParam(r, "username")
	follower, err := h.app.Dao.Account().FindByUsername(ctx, followerName)
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}
	if follower == nil {
		httperror.Error(w, http.StatusNotFound)
		return
	}

	relation := new(object.RelationWith)
	relationRepo := h.app.Dao.Relation()
	err = relationRepo.Unfollow(ctx, requester.ID, follower.ID)
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}

	relation.Following, err = relationRepo.IsFollowing(ctx, requester.ID, follower.ID)
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}

	relation.FollowedBy, err = relationRepo.IsFollowing(ctx, follower.ID, requester.ID)
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
