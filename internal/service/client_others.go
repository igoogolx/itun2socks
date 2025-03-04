//go:build !windows

package service

import (
	"fmt"
)

var err = fmt.Errorf("not implemented yet")

func Run() error {
	return err
}

func Install() error {
	return err
}

func Uninstall() error {
	return err
}

func Restart() error {
	return err
}

func Interactive() bool {
	return false
}
