package object

type (
	StatusID = int64

	Status struct {
		ID       StatusID `json:"id" db:"id"`
		Account  Account  `json:"account"`
		Content  string   `json:"content" db:"content"`
		CreateAt DateTime `json:"create_at,omitempty" db:"create_at"`
		Media    Media    `json:"media_attachments"`
	}

	MediaID = int64
	Media   struct {
		ID MediaID `json:"id"`
	}
)
