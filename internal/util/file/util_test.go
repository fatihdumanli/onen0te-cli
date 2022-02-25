package file

import (
	"os"
	"strconv"
	"testing"
)

func Test_Exists(t *testing.T) {
	os.Create("/tmp/temp.file")
	defer os.Remove("/tmp/temp.file")

	if !Exists("/tmp/temp.file") {
		t.Error("in Test_Exists: expected true, got false")
	}
}

//Reads the contents of a file and returns the contents
func Test_ReadString(t *testing.T) {
}

func Test_HumanizeSize(t *testing.T) {
	data := []struct {
		numOfBytes int
		output     string
	}{
		{1024, "1 kB"},
		{5000, "5 kB"},
		{10240, "10 kB"},
	}

	for _, d := range data {
		t.Run(strconv.Itoa(d.numOfBytes), func(t *testing.T) {
			res := HumanizeSize(d.numOfBytes)

			if res != d.output {
				t.Errorf("expected %s, got %s", d.output, res)
			}
		})
	}

}
