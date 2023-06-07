package list

import (
	"github.com/Dreamacro/clash/component/geodata/router"
	"net"
)

type Lister struct {
	Items  []string
	Mather func(s, i string) bool
}

func (l Lister) Has(i string) bool {
	for _, item := range l.Items {
		if l.Mather(item, i) {
			return true
		}
	}
	return false
}

func New(items []string, matcher func(s, i string) bool) *Lister {
	return &Lister{
		Items:  items,
		Mather: matcher,
	}
}

type IpLister struct {
	Items []*router.GeoIPMatcher
}

func (l IpLister) Has(i string) bool {
	ip := net.ParseIP(i)
	if ip == nil {
		return false
	}
	for _, item := range l.Items {
		if item.Match(ip) {
			return true
		}
	}
	return false
}

type DomainLister struct {
	Items []*router.DomainMatcher
}

func (l DomainLister) Has(i string) bool {
	for _, item := range l.Items {
		if item.ApplyDomain(i) {
			return true
		}
	}
	return false
}
