package tun

import (
	"errors"
	"net"
)

func ParseAddresses(localAddr, netMask, nextHop string) (*net.IPNet, *net.IPNet, error) {
	local := net.ParseIP(localAddr)
	if local == nil {
		return nil, nil, errors.New("invalid local IP address")
	}

	mask := net.ParseIP(netMask)
	if mask == nil {
		return nil, nil, errors.New("invalid local IP mask")
	}

	gw := net.ParseIP(nextHop)
	if gw == nil {
		return nil, nil, errors.New("invalid gateway IP address")
	}

	loc := &net.IPNet{
		IP:   local.To4(),
		Mask: net.IPMask(mask.To4()),
	}
	rem := &net.IPNet{
		IP:   gw.To4(),
		Mask: net.CIDRMask(32, 32),
	}

	return loc, rem, nil
}
