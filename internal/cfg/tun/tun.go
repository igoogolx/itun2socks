package tun

import (
	"github.com/igoogolx/itun2socks/internal/constants"
	sTun "github.com/sagernet/sing-tun"
	"net/netip"
)

type Config struct {
	Name    string
	Gateway netip.Prefix
	Mtu     int
}

func New() (*Config, error) {
	gw, err := ParseAddresses(constants.TunGateway)
	if err != nil {
		return nil, err
	}
	tunInterfaceName := sTun.CalculateInterfaceName(constants.TunName)
	return &Config{
		Name:    tunInterfaceName,
		Mtu:     constants.TunMtu,
		Gateway: gw,
	}, nil
}
