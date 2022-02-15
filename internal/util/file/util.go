package file

import (
	"fmt"
	"io"
	"os"

	"github.com/pkg/errors"
)

func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
func ReadString(path string) (string, error) {

	var content string
	f, err := os.Open(path)
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("couldn't open the file %s", path))
	}

	bytes, err := io.ReadAll(f)
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("couldn't read the file %s", path))
	}
	content = string(bytes)
	return content, nil
}
