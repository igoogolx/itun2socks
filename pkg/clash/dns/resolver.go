package dns

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"strings"
	"time"

	"github.com/Dreamacro/clash/common/cache"
	"github.com/Dreamacro/clash/component/fakeip"
	"github.com/Dreamacro/clash/component/resolver"
	"github.com/Dreamacro/clash/component/trie"
	C "github.com/Dreamacro/clash/constant"

	D "github.com/miekg/dns"
	"github.com/samber/lo"
	"golang.org/x/sync/singleflight"
)

type dnsClient interface {
	GetServers() []string
	Exchange(m *D.Msg) (msg *D.Msg, err error)
	ExchangeContext(ctx context.Context, m *D.Msg) (msg *D.Msg, err error)
}

type result struct {
	Msg   *D.Msg
	Error error
}

type Resolver struct {
	ipv6          bool
	hosts         *trie.DomainTrie
	main          []dnsClient
	fallback      []dnsClient
	group         singleflight.Group
	lruCache      *cache.LruCache
	policy        *trie.DomainTrie
	searchDomains []string
	disableCache  bool
}

func (r *Resolver) GetServers() []string {
	var servers []string
	for _, c := range r.main {
		servers = append(servers, c.GetServers()...)
	}
	return servers
}

// LookupIP request with TypeA and TypeAAAA, priority return TypeA
func (r *Resolver) LookupIP(ctx context.Context, host string) (ip []net.IP, err error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	ch := make(chan []net.IP, 1)

	go func() {
		defer close(ch)
		ip, err := r.lookupIP(ctx, host, D.TypeAAAA)
		if err != nil {
			return
		}
		ch <- ip
	}()

	ip, err = r.lookupIP(ctx, host, D.TypeA)
	if err == nil {
		return
	}

	ip, open := <-ch
	if !open {
		return nil, resolver.ErrIPNotFound
	}

	return ip, nil
}

// ResolveIP request with TypeA and TypeAAAA, priority return TypeA
func (r *Resolver) ResolveIP(host string) (ip net.IP, err error) {
	ips, err := r.LookupIP(context.Background(), host)
	if err != nil {
		return nil, err
	} else if len(ips) == 0 {
		return nil, fmt.Errorf("%w: %s", resolver.ErrIPNotFound, host)
	}
	return ips[rand.Intn(len(ips))], nil
}

// LookupIPv4 request with TypeA
func (r *Resolver) LookupIPv4(ctx context.Context, host string) ([]net.IP, error) {
	return r.lookupIP(ctx, host, D.TypeA)
}

// ResolveIPv4 request with TypeA
func (r *Resolver) ResolveIPv4(host string) (ip net.IP, err error) {
	ips, err := r.lookupIP(context.Background(), host, D.TypeA)
	if err != nil {
		return nil, err
	} else if len(ips) == 0 {
		return nil, fmt.Errorf("%w: %s", resolver.ErrIPNotFound, host)
	}
	return ips[rand.Intn(len(ips))], nil
}

// LookupIPv6 request with TypeAAAA
func (r *Resolver) LookupIPv6(ctx context.Context, host string) ([]net.IP, error) {
	return r.lookupIP(ctx, host, D.TypeAAAA)
}

// ResolveIPv6 request with TypeAAAA
func (r *Resolver) ResolveIPv6(host string) (ip net.IP, err error) {
	ips, err := r.lookupIP(context.Background(), host, D.TypeAAAA)
	if err != nil {
		return nil, err
	} else if len(ips) == 0 {
		return nil, fmt.Errorf("%w: %s", resolver.ErrIPNotFound, host)
	}
	return ips[rand.Intn(len(ips))], nil
}

// Exchange a batch of dns request, and it use cache
func (r *Resolver) Exchange(m *D.Msg) (msg *D.Msg, err error) {
	return r.ExchangeContext(context.Background(), m)
}

// ExchangeContext a batch of dns request with context.Context, and it use cache
func (r *Resolver) ExchangeContext(ctx context.Context, m *D.Msg) (msg *D.Msg, err error) {
	if len(m.Question) == 0 {
		return nil, errors.New("should have one question at least")
	}

	q := m.Question[0]
	cacheItem, expireTime, hit := r.lruCache.GetWithExpire(q.String())
	cachedMsg, ok := cacheItem.(*D.Msg)
	if hit && ok {
		now := time.Now()
		msg = copyMsgFromCache(m, cachedMsg)
		if expireTime.Before(now) {
			setMsgTTL(msg, uint32(1)) // Continue fetch
			go func() {
				ctx, cancel := context.WithTimeout(context.Background(), resolver.DefaultDNSTimeout)
				r.exchangeWithoutCache(ctx, m)
				cancel()
			}()
		} else {
			// updating TTL by subtracting common delta time from each DNS record
			updateMsgTTL(msg, uint32(time.Until(expireTime).Seconds()))
		}
		return
	}
	return r.exchangeWithoutCache(ctx, m)
}

func (r *Resolver) exchange(ctx context.Context, m *D.Msg) (msg *D.Msg, err error) {
	q := m.Question[0]
	isIPReq := isIPRequest(q)
	if isIPReq {
		return r.ipExchange(ctx, m)
	}

	if matched := r.matchPolicy(m); len(matched) != 0 {
		return r.batchExchange(ctx, matched, m)
	}
	return r.batchExchange(ctx, r.main, m)
}

// ExchangeWithoutCache a batch of dns request, and it do NOT GET from cache
func (r *Resolver) exchangeWithoutCache(ctx context.Context, m *D.Msg) (*D.Msg, error) {
	if r.disableCache {
		return r.exchange(ctx, m)
	}
	q := m.Question[0]
	ret, err, shared := r.group.Do(q.String(), func() (result any, err error) {
		defer func() {
			if err != nil {
				return
			}

			msg := result.(*D.Msg)
			// OPT RRs MUST NOT be cached, forwarded, or stored in or loaded from master files.
			msg.Extra = lo.Filter(msg.Extra, func(rr D.RR, index int) bool {
				return rr.Header().Rrtype != D.TypeOPT
			})
			putMsgToCache(r.lruCache, q.String(), q, msg)
		}()
		return r.exchange(ctx, m)
	})
	if err != nil {
		return nil, err
	}
	res := ret.(*D.Msg)
	if shared {
		res = res.Copy()
	}
	return res, nil

}

func (r *Resolver) batchExchange(ctx context.Context, clients []dnsClient, m *D.Msg) (msg *D.Msg, err error) {
	ctx, cancel := context.WithTimeout(ctx, resolver.DefaultDNSTimeout)
	defer cancel()

	return batchExchange(ctx, clients, m)
}

func (r *Resolver) matchPolicy(m *D.Msg) []dnsClient {
	if r.policy == nil {
		return nil
	}

	domain := r.msgToDomain(m)
	if domain == "" {
		return nil
	}

	record := r.policy.Search(domain)
	if record == nil {
		return nil
	}

	return record.Data.([]dnsClient)
}

func (r *Resolver) ipExchange(ctx context.Context, m *D.Msg) (msg *D.Msg, err error) {
	if matched := r.matchPolicy(m); len(matched) != 0 {
		res := <-r.asyncExchange(ctx, matched, m)
		return res.Msg, res.Error
	}

	msgCh := r.asyncExchange(ctx, r.main, m)

	if r.fallback == nil { // directly return if no fallback servers are available
		res := <-msgCh
		msg, err = res.Msg, res.Error
		return
	}

	fallbackMsg := r.asyncExchange(ctx, r.fallback, m)
	res := <-msgCh
	if res.Error == nil {
		if ips := msgToIP(res.Msg); len(ips) != 0 {
			msg = res.Msg // no need to wait for fallback result
			err = res.Error
			return msg, err
		}
	}

	res = <-fallbackMsg
	msg, err = res.Msg, res.Error
	return
}

func (r *Resolver) lookupIP(ctx context.Context, host string, dnsType uint16) ([]net.IP, error) {
	ip := net.ParseIP(host)
	if ip != nil {
		ip4 := ip.To4()
		isIPv4 := ip4 != nil
		if dnsType == D.TypeAAAA && !isIPv4 {
			return []net.IP{ip}, nil
		} else if dnsType == D.TypeA && isIPv4 {
			return []net.IP{ip4}, nil
		} else {
			return nil, resolver.ErrIPVersion
		}
	}

	query := &D.Msg{}
	query.SetQuestion(D.Fqdn(host), dnsType)

	msg, err := r.ExchangeContext(ctx, query)
	if err != nil {
		return nil, err
	}

	ips := msgToIP(msg)
	if len(ips) != 0 {
		return ips, nil
	} else if len(r.searchDomains) == 0 {
		return nil, resolver.ErrIPNotFound
	}

	// query provided search domains serially
	for _, domain := range r.searchDomains {
		q := &D.Msg{}
		q.SetQuestion(D.Fqdn(fmt.Sprintf("%s.%s", host, domain)), dnsType)
		msg, err := r.ExchangeContext(ctx, q)
		if err != nil {
			return nil, err
		}
		ips := msgToIP(msg)
		if len(ips) != 0 {
			return ips, nil
		}
	}

	return nil, resolver.ErrIPNotFound
}

func (r *Resolver) msgToDomain(msg *D.Msg) string {
	if len(msg.Question) > 0 {
		return strings.TrimRight(msg.Question[0].Name, ".")
	}

	return ""
}

func (r *Resolver) asyncExchange(ctx context.Context, client []dnsClient, msg *D.Msg) <-chan *result {
	ch := make(chan *result, 1)
	go func() {
		res, err := r.batchExchange(ctx, client, msg)
		ch <- &result{Msg: res, Error: err}
	}()
	return ch
}

type NameServer struct {
	Net       string
	Addr      string
	Interface string
}

type FallbackFilter struct {
	GeoIP     bool
	GeoIPCode string
	IPCIDR    []*net.IPNet
	Domain    []string
}

type Config struct {
	Main, Fallback []NameServer
	Default        []NameServer
	IPv6           bool
	EnhancedMode   C.DNSMode
	FallbackFilter FallbackFilter
	Pool           *fakeip.Pool
	Hosts          *trie.DomainTrie
	Policy         map[string]NameServer
	SearchDomains  []string
	DisableCache   bool
	GetDialer      func() (C.Proxy, error)
}

func NewResolver(config Config) *Resolver {

	r := &Resolver{
		ipv6:          config.IPv6,
		main:          transform(config.Main, config.GetDialer, config.Pool),
		lruCache:      cache.New(cache.WithSize(4096), cache.WithStale(true)),
		hosts:         config.Hosts,
		searchDomains: config.SearchDomains,
		disableCache:  config.DisableCache,
	}

	if len(config.Fallback) != 0 {
		r.fallback = transform(config.Fallback, config.GetDialer, config.Pool)
	}

	if len(config.Policy) != 0 {
		r.policy = trie.New()
		for domain, nameserver := range config.Policy {
			r.policy.Insert(domain, transform([]NameServer{nameserver}, config.GetDialer, config.Pool))
		}
	}

	return r
}
