package file

import (
	"fmt"
	"io"
	"math"
	"os"
	"strconv"

	"github.com/pkg/errors"
)

//Returns true if the file exists in the given path
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

func HumanizeSize(numOfBytes int) string {
	if numOfBytes < 1024 {
		return strconv.Itoa(numOfBytes) + "bytes"
	}
	x := float64(numOfBytes) / float64(1024)
	return strconv.Itoa(int(math.Ceil(x))) + " kB"
}
