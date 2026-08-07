package resolver

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"net/netip"
	"strings"
	"time"

	"github.com/igoogolx/itun2socks/pkg/clash/component/trie"

	"github.com/miekg/dns"
)

var (
	// DefaultResolver aim to resolve ip
	DefaultResolver Resolver

	// DisableIPv6 means don't resolve ipv6 host
	// default value is true
	DisableIPv6 = true

	// DefaultHosts aim to resolve hosts
	DefaultHosts = trie.New()

	// DefaultDNSTimeout defined the default dns request timeout
	DefaultDNSTimeout = time.Second * 5
)

var (
	ErrIPNotFound   = errors.New("couldn't find ip")
	ErrIPVersion    = errors.New("ip version error")
	ErrIPv6Disabled = errors.New("ipv6 disabled")
)

type Resolver interface {
	GetServers() []string
	LookupIP(ctx context.Context, host string) ([]netip.Addr, error)
	LookupIPv4(ctx context.Context, host string) ([]netip.Addr, error)
	LookupIPv6(ctx context.Context, host string) ([]netip.Addr, error)
	ResolveIP(host string) (ip netip.Addr, err error)
	ResolveIPv4(host string) (ip netip.Addr, err error)
	ResolveIPv6(host string) (ip netip.Addr, err error)
	ExchangeContext(ctx context.Context, m *dns.Msg) (msg *dns.Msg, err error)
}

// LookupIPv4 with a host, return ipv4 list
func LookupIPv4(ctx context.Context, host string) ([]netip.Addr, error) {
	if node := DefaultHosts.Search(host); node != nil {
		if ip, ok := node.Data.(netip.Addr); ok && ip.Is4() {
			return []netip.Addr{ip}, nil
		}
	}

	ip, err := netip.ParseAddr(host)
	if err == nil {
		if !strings.Contains(host, ":") {
			return []netip.Addr{ip}, nil
		}
		return nil, ErrIPVersion
	}

	if DefaultResolver != nil {
		return DefaultResolver.LookupIPv4(ctx, host)
	}

	ctx, cancel := context.WithTimeout(context.Background(), DefaultDNSTimeout)
	defer cancel()
	ipAddrs, err := net.DefaultResolver.LookupNetIP(ctx, "ip4", host)
	if err != nil {
		return nil, err
	} else if len(ipAddrs) == 0 {
		return nil, ErrIPNotFound
	}

	return ipAddrs, nil
}

// ResolveIPv4 with a host, return ipv4
func ResolveIPv4(host string) (netip.Addr, error) {
	ips, err := LookupIPv4(context.Background(), host)
	if err != nil {
		return netip.Addr{}, err
	} else if len(ips) == 0 {
		return netip.Addr{}, fmt.Errorf("%w: %s", ErrIPNotFound, host)
	}
	return ips[rand.Intn(len(ips))], nil
}

// LookupIPv6 with a host, return ipv6 list
func LookupIPv6(ctx context.Context, host string) ([]netip.Addr, error) {
	if DisableIPv6 {
		return nil, ErrIPv6Disabled
	}

	if node := DefaultHosts.Search(host); node != nil {
		if ip, ok := node.Data.(netip.Addr); ok && ip.Is6() {
			return []netip.Addr{ip}, nil
		}
	}

	ip, err := netip.ParseAddr(host)
	if err == nil {
		if strings.Contains(host, ":") {
			return []netip.Addr{ip}, nil
		}
		return nil, ErrIPVersion
	}

	if DefaultResolver != nil {
		return DefaultResolver.LookupIPv6(ctx, host)
	}

	ctx, cancel := context.WithTimeout(context.Background(), DefaultDNSTimeout)
	defer cancel()
	ipAddrs, err := net.DefaultResolver.LookupNetIP(ctx, "ip6", host)
	if err != nil {
		return nil, err
	} else if len(ipAddrs) == 0 {
		return nil, ErrIPNotFound
	}

	return ipAddrs, nil
}

// ResolveIPv6 with a host, return ipv6
func ResolveIPv6(host string) (netip.Addr, error) {
	ips, err := LookupIPv6(context.Background(), host)
	if err != nil {
		return netip.Addr{}, err
	} else if len(ips) == 0 {
		return netip.Addr{}, fmt.Errorf("%w: %s", ErrIPNotFound, host)
	}
	return ips[rand.Intn(len(ips))], nil
}

// LookupIPWithResolver same as ResolveIP, but with a resolver
func LookupIPWithResolver(ctx context.Context, host string, r Resolver) ([]netip.Addr, error) {
	if node := DefaultHosts.Search(host); node != nil {
		return []netip.Addr{node.Data.(netip.Addr)}, nil
	}

	if r != nil {
		if DisableIPv6 {
			return r.LookupIPv4(ctx, host)
		}
		return r.LookupIP(ctx, host)
	} else if DisableIPv6 {
		return LookupIPv4(ctx, host)
	}

	ip, err := netip.ParseAddr(host)
	if err == nil {
		return []netip.Addr{ip}, nil
	}

	ips, err := net.DefaultResolver.LookupNetIP(ctx, "ip", host)
	if err != nil {
		return nil, err
	} else if len(ips) == 0 {
		return nil, ErrIPNotFound
	}

	return ips, nil
}

// LookupIP with context and a host, return ip
func LookupIP(ctx context.Context, host string) ([]netip.Addr, error) {
	return LookupIPWithResolver(ctx, host, DefaultResolver)
}

// ResolveIP with a host, return ip
func ResolveIP(host string) (netip.Addr, error) {
	ips, err := LookupIP(context.Background(), host)
	if err != nil {
		return netip.Addr{}, err
	} else if len(ips) == 0 {
		return netip.Addr{}, fmt.Errorf("%w: %s", ErrIPNotFound, host)
	}
	return ips[rand.Intn(len(ips))], nil
}
