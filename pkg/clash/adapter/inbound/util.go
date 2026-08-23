package inbound

import (
	"fmt"
	"net"
	"net/http"
	"net/netip"
	"strconv"
	"strings"

	"github.com/igoogolx/itun2socks/pkg/clash/common/util"
	C "github.com/igoogolx/itun2socks/pkg/clash/constant"
	"github.com/igoogolx/itun2socks/pkg/clash/transport/socks5"
)

func parseSocksAddr(target socks5.Addr) *C.Metadata {
	metadata := &C.Metadata{}

	switch target[0] {
	case socks5.AtypDomainName:
		// trim for FQDN
		metadata.Host = strings.TrimRight(string(target[2:2+target[1]]), ".")
		metadata.DstPort = C.Port((int(target[2+target[1]]) << 8) | int(target[2+target[1]+1]))
	case socks5.AtypIPv4:
		ip, ok := netip.AddrFromSlice(target[1 : 1+net.IPv4len])
		if ok {
			metadata.DstIP = ip
		}
		metadata.DstPort = C.Port((int(target[1+net.IPv4len]) << 8) | int(target[1+net.IPv4len+1]))
	case socks5.AtypIPv6:
		ip, ok := netip.AddrFromSlice(target[1 : 1+net.IPv6len])
		if ok {
			metadata.DstIP = ip
		}
		metadata.DstPort = C.Port((int(target[1+net.IPv6len]) << 8) | int(target[1+net.IPv6len+1]))
	}

	return metadata
}

func parseHTTPAddr(request *http.Request) *C.Metadata {
	host := request.URL.Hostname()
	port, _ := strconv.ParseUint(util.EmptyOr(request.URL.Port(), "80"), 10, 16)

	// trim FQDN (#737)
	host = strings.TrimRight(host, ".")

	metadata := &C.Metadata{
		NetWork: C.TCP,
		Host:    host,
		DstIP:   netip.Addr{},
		DstPort: C.Port(port),
	}

	if ip, err := netip.ParseAddr(host); err == nil {
		metadata.DstIP = ip
	}

	return metadata
}

func parseAddr(addr net.Addr) (netip.Addr, int, error) {
	switch a := addr.(type) {
	case *net.TCPAddr:
		netipAddr, ok := netip.AddrFromSlice(a.IP)
		if !ok {
			return netip.Addr{}, 0, fmt.Errorf("invalid tcp address %s", addr.String())
		}
		return netipAddr, a.Port, nil
	case *net.UDPAddr:

		netipAddr, ok := netip.AddrFromSlice(a.IP)
		if !ok {
			return netip.Addr{}, 0, fmt.Errorf("invalid udp address %s", addr.String())
		}
		return netipAddr, a.Port, nil
	default:
		return netip.Addr{}, 0, fmt.Errorf("unknown address type %s", addr.String())
	}
}
