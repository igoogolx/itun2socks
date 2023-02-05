//go:build windows

package dns

import (
	"fmt"
	"os/exec"
)

func FlushSysCaches() error {
	out, err := exec.Command("ipconfig", "/flushdns").CombinedOutput()
	if err != nil {
		return fmt.Errorf("%v (output: %s)", err, out)
	}
	return nil
}
