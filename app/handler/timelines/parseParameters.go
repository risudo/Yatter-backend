package timelines

import (
	"errors"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"yatter-backend-go/app/domain/object"
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
