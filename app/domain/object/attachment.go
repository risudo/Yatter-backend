package object

type (
	AttachmentID = int64

	Attachment struct {
		// ID of the attachment
		ID AttachmentID `json:"id"`

		// One of: "image", "video", "gifv", "unknown"
		MediaType string `json:"type" db:"type"`

		// URL of the image
		URL string `json:"url"`

		// A description of the image for the visually impaired (maximum 420 characters), or null if none provided
		Description *string `json:"desctiption"`
	}
)
