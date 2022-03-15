package object

type (
	FollowID = int64

	Relation struct {
		ID FollowID `db:"id"`

		FollowingID AccountID `db:"following_id"`

		FollowerID AccountID `db:"follower_id"`
	}
)
