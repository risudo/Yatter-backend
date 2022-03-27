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

var errOutOfRange = errors.New("limit is out of range")

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

	if value := r.FormValue("max_id"); value != "" {
		p.MaxID, err = strconv.ParseInt(value, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("parseParameters: %w", err)
		}
	}
	if value := r.FormValue("since_id"); value != "" {
		p.SinceID, err = strconv.ParseInt(value, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("parseParameters: %w", err)
		}
	}
	if value := r.FormValue("limit"); value != "" {
		p.Limit, err = strconv.Atoi(value)
		if err != nil {
			return nil, fmt.Errorf("parseParameters: %w", err)
		}
		if p.Limit > maxLimit || p.Limit < minLimit {
			return nil, errOutOfRange
		}
	}
	return p, nil
}

// Handler request for "GET /v1/timelines/public"
func (h *handler) Public(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	parameters, err := parseParameters(r)

	if err != nil {
		switch err {
		case errOutOfRange:
			httperror.Error(w, http.StatusBadRequest)
			return
		default:
			httperror.InternalServerError(w, err)
			return
		}
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
