//go:build windows
// +build windows

package is_elevated

import (
	"golang.org/x/sys/windows"
)

func Get() bool {
	var isAdmin bool
	var hToken = new(windows.Token)
	err := windows.OpenProcessToken(windows.CurrentProcess(), windows.TOKEN_QUERY, hToken)
	if err == nil {
		return hToken.IsElevated()
	}
	return isAdmin
}
