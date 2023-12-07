package tun

import (
	"net/netip"
)

func ParseAddresses(nextHop string) (netip.Prefix, error) {
	return netip.MustParsePrefix(nextHop), nil

}
