// +build linux solaris freebsd netbsd openbsd dragonfly

package clip

import (
	"os/exec"
	"regexp"

	"github.com/pkg/errors"
)

func getClipboard() ([]byte, error) {

	if _, err := exec.LookPath("xclip"); err != nil {
		return nil, errors.New("Could not found xclip. Please install it.")
	}

	c := exec.Command("xclip", "-o", "-selection", "c", "-t", "TARGETS")

	out, err := c.CombinedOutput()

	if err != nil {
		return nil, errors.New("Check TARGETS: xclip return nothing.")
	}

	r := regexp.MustCompile("image/png")

	if len(r.Find(out)) == 0 {
		return nil, errors.New("Clipboard can't containt picture.")
	}

	c = exec.Command("xclip", "-o", "-selection", "c", "-t", "image/png")

	out, err = c.Output()
	if err != nil {
		return nil, errors.New("Get picture: xclip return nothing.")
	}

	return out, nil
}
