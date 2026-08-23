package fakeip

import (
	"errors"
	"net/netip"
	"strings"
	"sync"

	"github.com/igoogolx/itun2socks/pkg/clash/common/cache"
	"github.com/igoogolx/itun2socks/pkg/clash/component/trie"
	"go4.org/netipx"
)

type store interface {
	GetByHost(host string) (netip.Addr, bool)
	PutByHost(host string, ip netip.Addr)
	GetByIP(ip netip.Addr) (string, bool)
	PutByIP(ip netip.Addr, host string)
	DelByIP(ip netip.Addr)
	Exist(ip netip.Addr) bool
	CloneTo(store)
}

// Pool is an implementation about fake ip generator without storage
type Pool struct {
	max     uint32
	min     uint32
	gateway uint32
	offset  uint32
	mux     sync.Mutex
	host    *trie.DomainTrie
	ipnet   *netip.Prefix
	store   store
}

// Lookup return a fake ip with host
func (p *Pool) Lookup(host string) netip.Addr {
	p.mux.Lock()
	defer p.mux.Unlock()

	// RFC4343: DNS Case Insensitive, we SHOULD return result with all cases.
	host = strings.ToLower(host)
	if ip, exist := p.store.GetByHost(host); exist {
		return ip
	}

	ip := p.get(host)
	p.store.PutByHost(host, ip)
	return ip
}

// LookBack return host with the fake ip
func (p *Pool) LookBack(ip netip.Addr) (string, bool) {
	p.mux.Lock()
	defer p.mux.Unlock()

	if !ip.Is4() {
		return "", false
	}

	return p.store.GetByIP(ip)
}

// ShouldSkipped return if domain should be skipped
func (p *Pool) ShouldSkipped(domain string) bool {
	if p.host == nil {
		return false
	}
	return p.host.Search(domain) != nil
}

// Exist returns if given ip exists in fake-ip pool
func (p *Pool) Exist(ip netip.Addr) bool {
	p.mux.Lock()
	defer p.mux.Unlock()

	if !ip.Is4() {
		return false
	}

	return p.store.Exist(ip)
}

// Gateway return gateway ip
func (p *Pool) Gateway() netip.Addr {
	return uintToIP(p.gateway)
}

// IPNet return raw ipnet
func (p *Pool) IPNet() *netip.Prefix {
	return p.ipnet
}

// CloneFrom clone cache from old pool
func (p *Pool) CloneFrom(o *Pool) {
	o.store.CloneTo(p.store)
}

func (p *Pool) get(host string) netip.Addr {
	current := p.offset
	for {
		ip := uintToIP(p.min + p.offset)
		if !p.store.Exist(ip) {
			break
		}

		p.offset = (p.offset + 1) % (p.max - p.min)
		// Avoid infinite loops
		if p.offset == current {
			p.offset = (p.offset + 1) % (p.max - p.min)
			ip := uintToIP(p.min + p.offset)
			p.store.DelByIP(ip)
			break
		}
	}
	ip := uintToIP(p.min + p.offset)
	p.store.PutByIP(ip, host)
	return ip
}

func ipToUint(addr netip.Addr) uint32 {
	ip := addr.AsSlice()
	v := uint32(ip[0]) << 24
	v += uint32(ip[1]) << 16
	v += uint32(ip[2]) << 8
	v += uint32(ip[3])
	return v
}

func uintToIP(v uint32) netip.Addr {
	addr, ok := netip.AddrFromSlice([]byte{byte(v >> 24), byte(v >> 16), byte(v >> 8), byte(v)})
	if !ok {
		return netip.Addr{}
	}
	return addr
}

type Options struct {
	IPNet *netip.Prefix
	Host  *trie.DomainTrie

	// Size sets the maximum number of entries in memory
	// and does not work if Persistence is true
	Size int
}

// New return Pool instance
func New(options Options) (*Pool, error) {

	var (
		hostAddr = options.IPNet.Masked().Addr()
		gateway  = hostAddr.Next()
		first    = gateway.Next().Next().Next() // default start with 198.18.0.4
		last     = netipx.PrefixLastIP(*options.IPNet)
	)

	if !options.IPNet.IsValid() || !first.IsValid() || !first.Less(last) {
		return nil, errors.New("ipnet don't have valid ip")
	}

	minValue := ipToUint(first)
	pool := &Pool{
		min:     minValue,
		max:     ipToUint(last),
		gateway: minValue - 1,
		host:    options.Host,
		ipnet:   options.IPNet,
	}
	pool.store = &memoryStore{
		cache: cache.New(cache.WithSize(options.Size * 2)),
	}

	return pool, nil
}
