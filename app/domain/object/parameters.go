package object

type Parameters struct {
	// Only return statuses that have media attachments (public and tag timelines only)
	OnlyMedia bool

	// Get a list of followings with ID less than this value
	MaxID int64

	// Get a list of followings with ID greater than this value
	SinceID int64

	// Maximum number of followings to get (Default 40, Max 80)
	Limit int
}
