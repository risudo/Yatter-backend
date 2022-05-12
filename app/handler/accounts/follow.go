package accounts

import (
	"encoding/json"
	"fmt"
	"net/http"
	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/handler/auth"
	"yatter-backend-go/app/handler/httperror"

	"github.com/go-chi/chi"
)

// Handle request for "POST /v1/accounts/{username}/follow"
func (h *handler) Follow(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	login := auth.AccountOf(r)
	if login == nil {
		httperror.InternalServerError(w, fmt.Errorf("lost account"))
		return
	}

	target, err := h.app.Dao.Account().FindByUsername(ctx, chi.URLParam(r, "username"))
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}
	if target == nil {
		httperror.Error(w, http.StatusNotFound)
		return
	}

	// relationshipを作成
	relation := &object.RelationShip{
		ID: target.ID,
	}

	// フォローしているか
	relation.Following, err = h.app.Dao.Relation().IsFollowing(ctx, login.ID, target.ID)
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}
	// フォローしてなかったらフォローする
	if !relation.Following {
		if err = h.app.Dao.Relation().Follow(ctx, login.ID, target.ID); err != nil {
			httperror.InternalServerError(w, err)
			return
		}
		relation.Following = true
	}

	// フォローされているか
	relation.FollowedBy, err = h.app.Dao.Relation().IsFollowing(ctx, target.ID, login.ID)
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
