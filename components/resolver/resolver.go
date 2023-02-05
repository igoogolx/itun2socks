package resolver

import (
	"context"
	"fmt"
	D "github.com/miekg/dns"
	"net"
	"strings"
)

type Client interface {
	Exchange(m *D.Msg) (msg *D.Msg, err error)
	ExchangeContext(ctx context.Context, m *D.Msg) (msg *D.Msg, err error)
	Nameservers() []string
}

func NewClient(namesever string) (Client, error) {
	if strings.HasPrefix(namesever, "dhcp") {
		return newDhcpClient(), nil
	} else if strings.HasPrefix(namesever, "https") {
		return NewDoHClient(namesever), nil
	}
	ip := net.ParseIP(namesever)
	if ip == nil {
		return nil, fmt.Errorf("invalid dns sever: %v", namesever)
	}
	return NewResolver([]net.IP{ip}), nil
}
