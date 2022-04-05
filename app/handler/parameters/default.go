package parameters

import (
	"math"
	"yatter-backend-go/app/domain/object"
)

const defaultLimit = 40
const maxLimit = 80
const minLimit = 0

func Default() *object.Parameters {
	return &object.Parameters{
		MaxID:   math.MaxInt64,
		SinceID: 0,
		Limit:   defaultLimit,
	}
}
