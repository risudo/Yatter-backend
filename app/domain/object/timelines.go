package object

type (
	Timelines []Status

	Parameters struct {
		MaxID   int64
		SinceID int64
		Limit   int
	}
)
