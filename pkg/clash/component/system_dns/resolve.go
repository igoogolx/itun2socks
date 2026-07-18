package system_dns

import (
	"net/netip"

	"github.com/samber/lo"
)

func ResolverV4Servers(ifaceName string) ([]string, error) {

	servers, err := ResolveServers(ifaceName)

	if err != nil {
		return nil, err
	}

	v4Servers := lo.Filter(servers, func(s string, index int) bool {

		itemAddr, err := netip.ParseAddr(s)

		if err != nil {
			return false
		}

		return itemAddr.Is4()

	})

	return v4Servers, err

}
