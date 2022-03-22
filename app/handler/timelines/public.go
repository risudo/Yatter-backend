package timelines

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/handler/httperror"
)

func parseParameters(r *http.Request) (*object.Parameters, error) {
	const maxLimit = 80
	const minLimit = 0
	const defaultLimit = 40

	var err error
	p := &object.Parameters{
		MaxID:   math.MaxInt64,
		SinceID: 0,
		Limit:   defaultLimit,
	}

	p.MaxID, err = strconv.ParseInt(r.FormValue("max_id"), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("parseParameters: %w", err)
	}

	p.SinceID, err = strconv.ParseInt(r.FormValue("since_id"), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("parseParameters: %w", err)
	}

	p.Limit, err = strconv.Atoi(r.FormValue("limit"))
	if err != nil {
		return nil, fmt.Errorf("parseParameters: %w", err)
	}
	if p.Limit > maxLimit {
		return nil, errors.New("parseParameters: limit is too large")
	}
	if p.Limit < minLimit {
		p.Limit = maxLimit
	}
	return p, nil
}

// Handler request for "GET /v1/timelines/public"
func (h *handler) Public(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	parameters, err := parseParameters(r)

	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}

	timeline, err := h.app.Dao.Status().PublicTimeline(ctx, parameters)
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(timeline); err != nil {
		httperror.InternalServerError(w, err)
		return
	}
}
