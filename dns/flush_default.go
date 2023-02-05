//go:build !windows
// +build !windows

package dns

func FlushSysCaches() error {
	return nil
}
