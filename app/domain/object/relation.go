package object

type (
	FollowID = int64

	Relation struct {
		ID FollowID `db:"id"`

		FolloweeID AccountID `db:"followee_id"`

		FollowerID AccountID `db:"follower_id"`
	}
)
