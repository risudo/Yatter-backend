package accounts

import (
	"encoding/json"
	"fmt"
	"net/http"
	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/handler/auth"
	"yatter-backend-go/app/handler/httperror"
)

// Handler request for "GET /v1/accounts/Relationships"
func (h *handler) Relationships(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	login := auth.AccountOf(r)
	if login == nil {
		httperror.InternalServerError(w, fmt.Errorf("lost account"))
		return
	}

	targetName := r.FormValue("username")
	target, err := h.app.Dao.Account().FindByUsername(ctx, targetName)
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	} else if target == nil {
		httperror.Error(w, http.StatusNotFound)
		return
	}

	relation := &object.RelationShip{
		ID: target.ID,
	}
	// フォローしているか
	relation.Following, err = h.app.Dao.Relation().IsFollowing(ctx, login.ID, target.ID)
	if err != nil {
		httperror.InternalServerError(w, err)
		return
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
