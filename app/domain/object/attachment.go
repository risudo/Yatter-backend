package object

type (
	AttachmentID = int64

	Attachment struct {
		// ID of the attachment
		ID AccountID `json:"id"`

		// One of: "image", "video", "gifv", "unknown"
		Atype string `json:"type"`

		// URL of the image
		URL string `json:"url"`

		// A description of the image for the visually impaired (maximum 420 characters), or null if none provided
		Description *string `json:"desctiption"`
	}
)
