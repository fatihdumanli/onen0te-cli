package internal

import (
	"io"
	"os"
)

func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
func ReadFile(path string) (string, bool) {

	var content string
	f, err := os.Open(path)
	if err != nil {
		return content, false
	}

	bytes, err := io.ReadAll(f)
	if err != nil {
		return content, false
	}
	content = string(bytes)
	return content, true
}
