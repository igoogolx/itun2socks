package system_dns

import (
	"golang.zx2c4.com/wireguard/windows/tunnel/winipcfg"
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
