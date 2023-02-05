//go:build windows
// +build windows

package executor

import (
	tun "github.com/sagernet/sing-tun"
)

func init() {
	tun.TunnelType = "Itun2socks"
}
