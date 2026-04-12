package dns

import (
	"context"
	"github.com/Dreamacro/clash/component/fakeip"
	"github.com/Dreamacro/clash/component/resolver"
	D "github.com/miekg/dns"
	"strings"
)

type fakeIpClient struct {
	pool *fakeip.Pool
}

func (f *fakeIpClient) GetServers() []string {
	return []string{"Fake Ip"}
}

func (f *fakeIpClient) Exchange(m *D.Msg) (msg *D.Msg, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), resolver.DefaultDNSTimeout)
	defer cancel()

	return f.ExchangeContext(ctx, m)
}

func (f *fakeIpClient) ExchangeContext(ctx context.Context, m *D.Msg) (msg *D.Msg, err error) {

	q := m.Question[0]

	host := strings.TrimRight(q.Name, ".")

	switch q.Qtype {
	case D.TypeAAAA, D.TypeSVCB, D.TypeHTTPS:
		return handleMsgWithEmptyAnswer(m), nil
	}

	rr := &D.A{}
	rr.Hdr = D.RR_Header{Name: q.Name, Rrtype: D.TypeA, Class: D.ClassINET, Ttl: dnsDefaultTTL}
	ip := f.pool.Lookup(host)
	rr.A = ip
	msg = m.Copy()
	msg.Answer = []D.RR{rr}

	setMsgTTL(msg, 1)
	msg.SetRcode(m, D.RcodeSuccess)
	msg.Authoritative = true
	msg.RecursionAvailable = true

	return msg, nil
}

func newFakeIpClient(pool *fakeip.Pool) *fakeIpClient {
	newClient := &fakeIpClient{
		pool,
	}
	return newClient
}
