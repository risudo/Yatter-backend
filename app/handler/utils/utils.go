package utils

import (
	"path/filepath"
	"time"
)

func CreateURL(filename string) string {
	return "attachments/" + time.Now().Format(time.RFC3339Nano) + filepath.Ext(filename)
}
