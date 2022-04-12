package files

import (
	"os"
	"path/filepath"
	"time"
)

const attachmentDir = "attachments/"

func CreateURL(filename string) string {
	return attachmentDir + time.Now().Format(time.RFC3339Nano) + filepath.Ext(filename)
}

func MightCreateAttachmentDir() {
	if f, err := os.Stat(attachmentDir); os.IsNotExist(err) || !f.IsDir() {
		os.Mkdir(attachmentDir, 0777)
	}
}
