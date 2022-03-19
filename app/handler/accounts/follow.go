package accounts

import (
	"encoding/json"
	"net/http"
	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/handler/auth"
	"yatter-backend-go/app/handler/httperror"

	"github.com/go-chi/chi"
)

// Handle request for "POST /v1/accounts/{username}/follow"
func (h *handler) Follow(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	following := auth.AccountOf(r)
	// TODO: チェックする必要ある?
	if following == nil {
		httperror.InternalServerError(w, nil) //
		return
	}

	followerName := chi.URLParam(r, "username")
	follower, err := h.app.Dao.Account().FindByUsername(ctx, followerName)
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}
	if follower == nil {
		httperror.Error(w, 404)
		return
	}

	// TODO: Relationshipsと被ってるからまとめたい
	relation := new(object.RelationWith)
	repo := h.app.Dao.Relation()
	relation.Following, err = repo.IsFollowing(ctx, following.ID, follower.ID)
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}
	if !relation.Following {
		err = repo.Follow(ctx, following.ID, follower.ID)
		if err != nil {
			httperror.InternalServerError(w, err)
			return
		}
		relation.Following = true
	}

	relation.FollowedBy, err = repo.IsFollowing(ctx, follower.ID, following.ID)
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
