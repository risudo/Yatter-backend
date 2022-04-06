package object

type (
	FollowID = int64

	// relationship with the target
	RelationShip struct {
		// Target account id
		ID AccountID `json:"id"`

		// Whether the user is currently following the account
		Following bool `json:"following"`

		// Whether the user is currently being followed by the account
		FollowedBy bool `json:"followed_by"`
	}
)
