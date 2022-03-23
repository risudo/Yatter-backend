package object

type (
	StatusID = int64

	Status struct {
		// ID of the status
		ID StatusID `json:"id" db:"id"`

		// account of the status
		Account *Account `json:"account"`

		// content of the status
		Content string `json:"content" db:"content"`

		// The time the account was created
		CreateAt DateTime `json:"create_at,omitempty" db:"create_at"`
	}
)
