//go:build !darwin && !linux && !freebsd

package redir

import (
	"errors"
	"net"

	"github.com/igoogolx/itun2socks/pkg/clash/transport/socks5"
)

func parserPacket(conn net.Conn) (socks5.Addr, error) {
	return nil, errors.New("system not support yet")
}
