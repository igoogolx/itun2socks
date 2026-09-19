package dns

import (
	"context"
	"net/netip"

	"github.com/igoogolx/itun2socks/pkg/log"
	metaResolver "github.com/metacubex/mihomo/component/resolver"
	"github.com/metacubex/mihomo/dns"
)

type resolversWithLog struct {
	*dns.Resolvers
}

func (r resolversWithLog) LookupIP(ctx context.Context, host string) ([]netip.Addr, error) {

	ips, err := r.Resolvers.LookupIP(ctx, host)

	log.Infoln(log.FormatLog(log.DnsPrefix, "Boost DNS, look up: %v, ips: %v, err:%v"), host, ips, err)

	return ips, err

}

func (r resolversWithLog) LookupIPv4(ctx context.Context, host string) ([]netip.Addr, error) {

	ips, err := r.Resolvers.LookupIPv4(ctx, host)

	log.Infoln(log.FormatLog(log.DnsPrefix, "Boost DNS, look up: %v, ips: %v, err:%v"), host, ips, err)

	return ips, err

}

func (r resolversWithLog) LookupIPv6(ctx context.Context, host string) ([]netip.Addr, error) {

	ips, err := r.Resolvers.LookupIPv6(ctx, host)

	log.Infoln(log.FormatLog(log.DnsPrefix, "Boost DNS, look up: %v, ips: %v, err:%v"), host, ips, err)

	return ips, err

}

func InitDnsForSysProxy() error {

	nameServers, err := dns.ParseNameServer([]string{"system"})
	if err != nil {
		return err
	}

	metaResolver.SystemResolver = dns.NewResolver(dns.Config{Main: nameServers, CacheMaxSize: 1})

	return nil

}

func InitDnsForTunProxy(sysResolvers dns.Resolvers) {
	ResetCache()
	r := resolversWithLog{&sysResolvers}
	metaResolver.DefaultResolver = r
	metaResolver.SystemResolver = r
}
