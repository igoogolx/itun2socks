package dns

import (
	"context"
	"github.com/Dreamacro/clash/component/resolver"
	"github.com/Dreamacro/clash/component/system_dns"
	C "github.com/Dreamacro/clash/constant"
	"github.com/Dreamacro/clash/log"
	"net"
	"net/netip"
	"sync"

	D "github.com/miekg/dns"
)

type systemClient struct {
	ifaceName string
	lock      sync.Mutex
	clients   []dnsClient
	getDialer func() (C.Proxy, error)
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
				log.Warnln("Batch exchange failed:", err)
			}
		}()
	}
	return mRes, err
}

func (s *systemClient) update() error {
	dns, err := system_dns.ResolveServers(s.ifaceName)
	if err != nil {
		return err
	} else {
		log.Infoln("System DNS resolve: %s\n", dns)
	}
	var res []dnsClient
	nameserver := make([]NameServer, 0, len(dns))
	for _, item := range dns {
		itemAddr, err := netip.ParseAddr(item)
		if err == nil && itemAddr.Is4() {
			nameserver = append(nameserver, NameServer{
				Addr:      net.JoinHostPort(item, "53"),
				Interface: s.ifaceName,
			})
		}

	}

	res = transform(nameserver, s.getDialer, nil)
	s.lock.Lock()
	s.clients = res
	s.lock.Unlock()
	return nil
}

func newSystemClient(ifaceName string, getDialer func() (C.Proxy, error)) *systemClient {
	newClient := &systemClient{ifaceName: ifaceName, getDialer: getDialer}
	err := newClient.update()
	if err != nil {
		log.Warnln("System DNS init failed:", err)
	}
	return newClient
}
