package system_dns

import (
	"github.com/metacubex/tailscale/util/winipcfg"
	"net"
)

func ResolveServers(ifaceName string) ([]string, error) {
	networkInterface, err := net.InterfaceByName(ifaceName)
	if err != nil {
		return nil, err
	}
	luid, err := winipcfg.LUIDFromIndex(uint32(networkInterface.Index))
	if err != nil {
		return nil, err
	}
	dnsServers, err := luid.DNS()
	if err != nil {
		return nil, err
	}
	servers := make([]string, 0, len(dnsServers))
	for _, server := range dnsServers {
		servers = append(servers, server.String())
	}
	return servers, nil
}
