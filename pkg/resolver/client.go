package resolver

import (
	"context"
	"github.com/igoogolx/itun2socks/internal/constants"
	"github.com/miekg/dns"
	"net"
	"sync"
	"time"
)

var (
	DnsDefaultTimeout = 5 * time.Second
)

type Resolver struct {
	nameservers []net.IP
	addresses   []string
}

type Result struct {
	msg *dns.Msg
	err error
}

func (r Resolver) Exchange(m *dns.Msg) (*dns.Msg, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DnsDefaultTimeout)
	defer cancel()
	return r.ExchangeContext(ctx, m)
}

func (r Resolver) ExchangeContext(ctx context.Context, m *dns.Msg) (*dns.Msg, error) {
	resultCh := make(chan Result, len(r.nameservers))
	wg := sync.WaitGroup{}
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	for _, address := range r.addresses {
		go func(a string) {
			wg.Add(1)
			var ct = &dns.Client{
				Net: "tcp",
			}
			msg, _, err := ct.ExchangeContext(ctx, m, a)
			resultCh <- Result{
				msg: msg,
				err: err,
			}
			wg.Done()
		}(address)
	}

	var rc Result
	for range r.nameservers {
		rc = <-resultCh
		if rc.err == nil {
			break
		}
	}
	go func() {
		wg.Wait()
		close(resultCh)
	}()
	return rc.msg, rc.err
}

func (r Resolver) Nameservers() []string {
	servers := make([]string, 0, len(r.nameservers))
	for _, n := range r.nameservers {
		servers = append(servers, n.String())
	}
	return servers
}

func NewResolver(nameservers []net.IP) Resolver {
	addresses := make([]string, 0, len(nameservers))
	for _, nameserver := range nameservers {
		addresses = append(addresses, nameserver.String()+":"+constants.DnsPort)
	}

	return Resolver{
		nameservers,
		addresses,
	}
}
