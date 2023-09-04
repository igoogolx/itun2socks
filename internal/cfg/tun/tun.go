package tun

import (
	"github.com/igoogolx/itun2socks/internal/constants"
	"net"
)

type Config struct {
	Name      string
	LocalAddr *net.IPNet
	Gateway   *net.IPNet
	Mtu       int
}

func New() (Config, error) {
	local, gw, err := ParseAddresses(
		constants.TunLocalAddr,
		constants.TunMask,
		constants.TunGateway,
	)
	if err != nil {
		return Config{}, err
	}
	return Config{
		Name:      constants.TunName,
		Mtu:       constants.TunMtu,
		LocalAddr: local,
		Gateway:   gw,
	}, nil
}
