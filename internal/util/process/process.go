package process

import (
	"log"
	"os/exec"
	"runtime"

	"github.com/pkg/errors"
)

func Start(url string) error {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		log.Fatal("your OS is not supported.")
	}

	if err != nil {
		return errors.Wrap(err, "couldn't start the web browser to authenticate the user")
	}
	return nil
}
