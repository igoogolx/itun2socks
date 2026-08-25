package dns

import (
	"context"
	"net"
	"sync"

	"github.com/igoogolx/itun2socks/pkg/clash/component/resolver"
	"github.com/igoogolx/itun2socks/pkg/clash/component/system_dns"
	"github.com/igoogolx/itun2socks/pkg/clash/log"
	metaC "github.com/metacubex/mihomo/constant"

	D "github.com/miekg/dns"
)

type systemClient struct {
	ifaceName string
	lock      sync.Mutex
	clients   []dnsClient
	getDialer func() (metaC.Proxy, error)
}

func (s *systemClient) GetServers() []string {
	var servers []string
	for _, c := range s.clients {
		servers = append(servers, c.GetServers()...)
	}
	return servers
}

func (s *systemClient) Exchange(m *D.Msg) (msg *D.Msg, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), resolver.DefaultDNSTimeout)
	defer cancel()

	return s.ExchangeContext(ctx, m)
}

func (s *systemClient) ExchangeContext(ctx context.Context, m *D.Msg) (msg *D.Msg, err error) {
	var clients = s.clients
	if len(clients) == 0 {
		err = s.update()
		if err != nil {
			return nil, err
		}
	}
	mRes, err := batchExchange(ctx, clients, m)
	if err != nil {
		go func() {
			err := s.update()
			if err != nil {
				log.Warnln("Batch exchange failed: %v", err)
			}
		}()
	}
	return mRes, err
}

func (s *systemClient) update() error {
	dns, err := system_dns.ResolverV4Servers(s.ifaceName)
	if err != nil {
		return err
	}

	log.Infoln("System DNS resolve: %s\n", dns)
	var res []dnsClient
	nameserver := make([]NameServer, 0, len(dns))
	for _, item := range dns {
		nameserver = append(nameserver, NameServer{
			Addr:      net.JoinHostPort(item, "53"),
			Interface: s.ifaceName,
		})

	}

	res = transform(nameserver, s.getDialer, nil)
	s.lock.Lock()
	s.clients = res
	s.lock.Unlock()
	return nil
}

func newSystemClient(ifaceName string, getDialer func() (metaC.Proxy, error)) *systemClient {
	newClient := &systemClient{ifaceName: ifaceName, getDialer: getDialer}
	err := newClient.update()
	if err != nil {
		log.Warnln("System DNS init failed: %v", err)
	}
	return newClient
}
