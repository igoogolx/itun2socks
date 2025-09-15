package rule_engine

import (
	"context"
	"fmt"
	"net"
	"net/netip"
	"strings"

	"github.com/Dreamacro/clash/component/trie"
	"github.com/igoogolx/itun2socks/internal/constants"
	D "github.com/miekg/dns"
)

var dnsDefaultTTL uint32 = 600

type DnsMap struct {
	RuleType   constants.RuleType `json:"ruleType"`
	Payload    string             `json:"payload"`
	domainTrie *trie.DomainTrie
	domain     string
	ip         string
	Policy     constants.Policy `json:"policy"`
}

func (d DnsMap) GetPolicy() constants.Policy {
	return d.Policy
}

func (d DnsMap) Valid() bool {
	if len(d.domain) == 0 {
		return false
	}
	if _, err := netip.ParseAddr(d.ip); err != nil {
		return false
	}
	return true
}

func (d DnsMap) Type() constants.RuleType {
	return constants.RuleDnsMap
}

func (d DnsMap) Match(value string) bool {
	record := d.domainTrie.Search(value)
	return record != nil
}

func (d DnsMap) Value() string {
	return d.Payload
}

func (d DnsMap) ExchangeContext(_ context.Context, r *D.Msg) (*D.Msg, error) {

	q := r.Question[0]
	ip := net.ParseIP(d.ip)
	if ip == nil {
		return nil, fmt.Errorf("invalid ip address: %s", d.ip)
	}

	msg := r.Copy()

	if v4 := ip.To4(); v4 != nil && q.Qtype == D.TypeA {
		rr := &D.A{}
		rr.Hdr = D.RR_Header{Name: q.Name, Rrtype: D.TypeA, Class: D.ClassINET, Ttl: dnsDefaultTTL}
		rr.A = v4

		msg.Answer = []D.RR{rr}
	} else if v6 := ip.To16(); v6 != nil && q.Qtype == D.TypeAAAA {
		rr := &D.AAAA{}
		rr.Hdr = D.RR_Header{Name: q.Name, Rrtype: D.TypeAAAA, Class: D.ClassINET, Ttl: dnsDefaultTTL}
		rr.AAAA = v6

		msg.Answer = []D.RR{rr}
	} else {
		return nil, fmt.Errorf("invalid ip type")
	}

	msg.SetRcode(r, D.RcodeSuccess)
	msg.Authoritative = true
	msg.RecursionAvailable = true
	return msg, nil
}

func NewDnsMapRule(payload string, policy constants.Policy) (*DnsMap, error) {
	chunks := trimArr(strings.Split(strings.TrimSpace(payload), ";"))
	if len(chunks) != 2 {
		return nil, fmt.Errorf("invald DNS-MAP rule")
	}

	domainTrie := trie.New()
	rule := &DnsMap{constants.RuleDnsMap, payload, domainTrie, chunks[0], chunks[1], policy}
	if !rule.Valid() {
		return nil, fmt.Errorf("invalid DNS-MAP rule")
	}
	err := domainTrie.Insert(rule.domain, rule.ip)
	if err != nil {
		return nil, err
	}
	return rule, nil
}
