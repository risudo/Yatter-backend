package accounts

import (
	"encoding/json"
	"log"
	"net/http"
	"yatter-backend-go/app/handler/auth"
	"yatter-backend-go/app/handler/httperror"
)

func (h *handler) Relationships(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	requester := auth.AccountOf(r)
	if requester == nil {
		log.Println("auth fail")//TODO
		return
	}

	username := r.FormValue("username")

	account, err := h.app.Dao.Account().FindByUsername(ctx, username)
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	} else if account == nil {
		// アカウントが見つからなかった場合
		return
	}

	relation, err := h.app.Dao.Relation().Relationship(ctx, account.ID, requester.ID)
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
