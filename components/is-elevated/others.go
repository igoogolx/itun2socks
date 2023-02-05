//go:build !windows
// +build !windows

package is_elevated

import (
	"os"
)

func Get() bool {
	return os.Geteuid() == 0
}
