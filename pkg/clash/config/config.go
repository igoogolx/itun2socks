package config

import (
	"fmt"
	"github.com/Dreamacro/clash/component/auth"
	"github.com/Dreamacro/clash/component/fakeip"
	"github.com/Dreamacro/clash/component/trie"
	C "github.com/Dreamacro/clash/constant"
	providerTypes "github.com/Dreamacro/clash/constant/provider"
	"github.com/Dreamacro/clash/dns"
	"github.com/Dreamacro/clash/log"
	T "github.com/Dreamacro/clash/tunnel"
	"net"
	"net/url"
	"strings"
)

// General config
type General struct {
	LegacyInbound
	Controller
	Authentication []string     `json:"authentication"`
	Mode           T.TunnelMode `json:"mode"`
	LogLevel       log.LogLevel `json:"log-level"`
	IPv6           bool         `json:"ipv6"`
	Interface      string       `json:"-"`
	RoutingMark    int          `json:"-"`
}

// Controller
type Controller struct {
	ExternalController string `json:"-"`
	ExternalUI         string `json:"-"`
	Secret             string `json:"-"`
}

type LegacyInbound struct {
	Port        int    `json:"port"`
	SocksPort   int    `json:"socks-port"`
	RedirPort   int    `json:"redir-port"`
	TProxyPort  int    `json:"tproxy-port"`
	MixedPort   int    `json:"mixed-port"`
	AllowLan    bool   `json:"allow-lan"`
	BindAddress string `json:"bind-address"`
}

// DNS config
type DNS struct {
	Enable            bool             `yaml:"enable"`
	IPv6              bool             `yaml:"ipv6"`
	NameServer        []dns.NameServer `yaml:"nameserver"`
	Fallback          []dns.NameServer `yaml:"fallback"`
	FallbackFilter    FallbackFilter   `yaml:"fallback-filter"`
	Listen            string           `yaml:"listen"`
	EnhancedMode      C.DNSMode        `yaml:"enhanced-mode"`
	DefaultNameserver []dns.NameServer `yaml:"default-nameserver"`
	FakeIPRange       *fakeip.Pool
	Hosts             *trie.DomainTrie
	NameServerPolicy  map[string]dns.NameServer
	SearchDomains     []string
}

// FallbackFilter config
type FallbackFilter struct {
	GeoIP     bool         `yaml:"geoip"`
	GeoIPCode string       `yaml:"geoip-code"`
	IPCIDR    []*net.IPNet `yaml:"ipcidr"`
	Domain    []string     `yaml:"domain"`
}

// Profile config
type Profile struct {
	StoreSelected bool `yaml:"store-selected"`
	StoreFakeIP   bool `yaml:"store-fake-ip"`
}

// Experimental config
type Experimental struct {
	UDPFallbackMatch bool `yaml:"udp-fallback-match"`
}

// Config is clash config manager
type Config struct {
	General      *General
	DNS          *DNS
	Experimental *Experimental
	Hosts        *trie.DomainTrie
	Profile      *Profile
	Inbounds     []C.Inbound
	Rules        []C.Rule
	Users        []auth.AuthUser
	Proxies      map[string]C.Proxy
	Providers    map[string]providerTypes.ProxyProvider
	Tunnels      []Tunnel
}

type RawDNS struct {
	Enable            bool              `yaml:"enable"`
	IPv6              *bool             `yaml:"ipv6"`
	UseHosts          bool              `yaml:"use-hosts"`
	NameServer        []string          `yaml:"nameserver"`
	Fallback          []string          `yaml:"fallback"`
	FallbackFilter    RawFallbackFilter `yaml:"fallback-filter"`
	Listen            string            `yaml:"listen"`
	EnhancedMode      C.DNSMode         `yaml:"enhanced-mode"`
	FakeIPRange       string            `yaml:"fake-ip-range"`
	FakeIPFilter      []string          `yaml:"fake-ip-filter"`
	DefaultNameserver []string          `yaml:"default-nameserver"`
	NameServerPolicy  map[string]string `yaml:"nameserver-policy"`
	SearchDomains     []string          `yaml:"search-domains"`
}

type RawFallbackFilter struct {
	GeoIP     bool     `yaml:"geoip"`
	GeoIPCode string   `yaml:"geoip-code"`
	IPCIDR    []string `yaml:"ipcidr"`
	Domain    []string `yaml:"domain"`
}

type tunnel struct {
	Network []string `yaml:"network"`
	Address string   `yaml:"address"`
	Target  string   `yaml:"target"`
	Proxy   string   `yaml:"proxy"`
}

type Tunnel tunnel

func ParseNameServer(servers []string) ([]dns.NameServer, error) {
	nameservers := []dns.NameServer{}

	for idx, server := range servers {
		// parse without scheme .e.g 8.8.8.8:53
		if !strings.Contains(server, "://") {
			server = "udp://" + server
		}
		u, err := url.Parse(server)
		if err != nil {
			return nil, fmt.Errorf("DNS NameServer[%d] format error: %s", idx, err.Error())
		}

		// parse with specific interface
		// .e.g 10.0.0.1#en0
		interfaceName := u.Fragment

		var addr, dnsNetType string
		switch u.Scheme {
		case "udp":
			addr, err = hostWithDefaultPort(u.Host, "53")
			dnsNetType = "" // UDP
		case "tcp":
			addr, err = hostWithDefaultPort(u.Host, "53")
			dnsNetType = "tcp" // TCP
		case "tls":
			addr, err = hostWithDefaultPort(u.Host, "853")
			dnsNetType = "tcp-tls" // DNS over TLS
		case "https":
			clearURL := url.URL{Scheme: "https", Host: u.Host, Path: u.Path, User: u.User}
			addr = clearURL.String()
			dnsNetType = "https" // DNS over HTTPS
		case "dhcp":
			addr = u.Host
			dnsNetType = "dhcp" // UDP from DHCP
		case "system":
			addr = u.Host
			dnsNetType = "system" // UDP from System
		case "fake-ip":
			addr = u.Host
			dnsNetType = "fake-ip" // Fake Ip
		default:
			return nil, fmt.Errorf("DNS NameServer[%d] unsupport scheme: %s", idx, u.Scheme)
		}

		if err != nil {
			return nil, fmt.Errorf("DNS NameServer[%d] format error: %s", idx, err.Error())
		}

		nameservers = append(
			nameservers,
			dns.NameServer{
				Net:       dnsNetType,
				Addr:      addr,
				Interface: interfaceName,
			},
		)
	}
	return nameservers, nil
}

func hostWithDefaultPort(host string, defPort string) (string, error) {
	if !strings.Contains(host, ":") {
		host += ":"
	}

	hostname, port, err := net.SplitHostPort(host)
	if err != nil {
		return "", err
	}

	if port == "" {
		port = defPort
	}

	return net.JoinHostPort(hostname, port), nil
}
