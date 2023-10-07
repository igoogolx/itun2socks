package tun

import (
	"github.com/igoogolx/itun2socks/internal/constants"
	sTun "github.com/sagernet/sing-tun"
	"net"
)

type Config struct {
	Name      string
	LocalAddr *net.IPNet
	Gateway   *net.IPNet
	Mtu       int
}

func New() (*Config, error) {
	local, gw, err := ParseAddresses(
		constants.TunLocalAddr,
		constants.TunMask,
		constants.TunGateway,
	)
	if err != nil {
		return nil, err
	}
	tunInterfaceName := sTun.CalculateInterfaceName(constants.TunName)
	return &Config{
		Name:      tunInterfaceName,
		Mtu:       constants.TunMtu,
		LocalAddr: local,
		Gateway:   gw,
	}, nil
}
