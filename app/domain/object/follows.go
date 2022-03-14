package object

type (
	FollowID = int64

	Follow struct {
		ID FollowID `db:"id"`

		AccountID AccountID `db:"account_id"`

		FollowID AccountID `db:"follow_account_id"`
	}

	Reration struct {
		Id FollowID `json:"id"`

		Folowing bool `json:"folloing"`

		FolowedBy bool `json:"followed_by"`
	}
)
