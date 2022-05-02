package parameters

import (
	"fmt"
	"net/http"
	"strconv"
	"yatter-backend-go/app/domain/object"
)

var ErrOutOfRange = fmt.Errorf("limit is out of range")
var ErrEmpty = fmt.Errorf("empty parameter")

func parseFormValue(r *http.Request, key string) (int64, error) {
	value := r.FormValue(key)
	if value == "" {
		return -1, ErrEmpty
	}
	intValue, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return -1, err
	}
	return intValue, nil
}

func parseLimitValue(limit int64) (int, error) {
	if limit > maxLimit || limit < minLimit {
		return -1, ErrOutOfRange
	}
	return int(limit), nil
}

func ParseLimit(r *http.Request) (int, error) {
	var intlimit int
	limit, err := parseFormValue(r, "limit")
	if err != nil && err != ErrEmpty {
		return -1, err
	} else if err != ErrEmpty {
		intlimit, err = parseLimitValue(limit)
		if err != nil {
			return -1, err
		}
	}
	if err == ErrEmpty {
		return DefaultLimit, nil
	}
	return intlimit, nil
}

func ParseAll(r *http.Request) (*object.Parameters, error) {
	var err error
	p := Default()

	only_media, err := parseFormValue(r, "only_media")
	if err != nil && err != ErrEmpty {
		return nil, err
	}
	if only_media != 0 && err != ErrEmpty {
		p.OnlyMedia = true
	}

	max_id, err := parseFormValue(r, "max_id")
	if err != nil && err != ErrEmpty {
		return nil, err
	} else if err != ErrEmpty {
		p.MaxID = max_id
	}

	since_id, err := parseFormValue(r, "since_id")
	if err != nil && err != ErrEmpty {
		return nil, err
	} else if err != ErrEmpty {
		p.SinceID = since_id
	}

	limit, err := ParseLimit(r)
	if err != nil {
		return nil, err
	}
	p.Limit = limit

	return p, nil
}
