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

// attachmentsディレクトリがなかったら作成
func MightCreateAttachmentDir() error {
	f, err := os.Stat(attachmentDir)
	if os.IsNotExist(err) || !f.IsDir() {
		os.Mkdir(attachmentDir, 0777)
	} else if err != nil {
		return err
	}

	return nil
}
