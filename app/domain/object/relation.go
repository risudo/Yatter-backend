package object

type (
	FollowID = int64

	// dbとのマッピング用構造体
	Relation struct {
		ID FollowID `db:"id"`

		FollowingID AccountID `db:"following_id"`

		FollowerID AccountID `db:"follower_id"`
	}

	// relationship with the target
	RelationWith struct {
		// Target account id
		ID AccountID `json:"id"`

		// Whether the user is currently following the account
		Following bool `json:"following"`

		// Whether the user is currently being followed by the account
		FollowedBy bool `json:"followed_by"`
	}
)
