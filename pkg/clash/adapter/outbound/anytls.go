package outbound

import (
	"context"
	"fmt"
	"net"
	"net/netip"
	"strconv"

	anytls "github.com/anytls/sing-anytls"
	"github.com/igoogolx/itun2socks/pkg/clash/component/dialer"
	C "github.com/igoogolx/itun2socks/pkg/clash/constant"
	"github.com/sagernet/sing/common/json/badoption"
	singMetadata "github.com/sagernet/sing/common/metadata"
)

type connContextKey string

type AnyTls struct {
	*Base
	client        *anytls.Client
	diaOutConnKey connContextKey
}

type AnyTlsOption struct {
	BasicOption
	Name                     string             `proxy:"name"`
	Server                   string             `proxy:"server"`
	Port                     int                `proxy:"port"`
	Password                 string             `proxy:"password,omitempty"`
	IdleSessionCheckInterval badoption.Duration `proxy:"idle-session-check-interval,omitempty"`
	IdleSessionTimeout       badoption.Duration `proxy:"idle-session-timeout,omitempty"`
	MinIdleSession           int                `proxy:"min-idle-session,omitempty"`
	UDP                      bool               `proxy:"udp,omitempty"`
}

func (at *AnyTls) StreamConn(c net.Conn, metadata *C.Metadata) (net.Conn, error) {

	if at.client == nil {
		return nil, fmt.Errorf("invalid AnyTls client")
	}

	dest := singMetadata.Socksaddr{
		Port: uint16(metadata.DstPort),
	}

	if len(metadata.Host) != 0 {
		dest.Fqdn = metadata.Host
	} else {

		addr, err := netip.ParseAddr(metadata.DstIP.String())

		if err != nil {

			return nil, fmt.Errorf("invalid addr for AnyTls client: %s, error: %s", addr, err)
		}
		dest.Addr = addr

	}

	return at.client.CreateProxy(context.WithValue(context.Background(), at.diaOutConnKey, c), dest)
}

func (at *AnyTls) DialContext(ctx context.Context, metadata *C.Metadata, opts ...dialer.Option) (_ C.Conn, err error) {

	c, err := dialer.DialContext(ctx, "tcp", at.addr, at.Base.DialOptions(opts...)...)
	if err != nil {
		return nil, fmt.Errorf("%s connect error: %w", at.addr, err)
	}
	tcpKeepAlive(c)

	//defer func(c net.Conn) {
	//	safeConnClose(c, err)
	//}(c)

	c, err = at.StreamConn(c, metadata)
	if err != nil {
		return nil, err
	}

	return NewConn(c, at), nil

}

func (at *AnyTls) dialOut(ctx context.Context) (net.Conn, error) {

	if v, ok := ctx.Value(at.diaOutConnKey).(net.Conn); ok {
		return v, nil
	}

	return nil, fmt.Errorf("no conn found")
}

func NewAnyTls(option AnyTlsOption) *AnyTls {

	ctx := context.Background()

	addr := net.JoinHostPort(option.Server, strconv.Itoa(option.Port))

	anyTlsClient := &AnyTls{

		Base: &Base{
			name:  option.Name,
			addr:  addr,
			tp:    C.AnyTls,
			udp:   option.UDP,
			iface: option.Interface,
			rmark: option.RoutingMark,
		},
		diaOutConnKey: connContextKey("connKey"),
	}

	client, err := anytls.NewClient(ctx, anytls.ClientConfig{
		Password:                 option.Password,
		IdleSessionCheckInterval: option.IdleSessionCheckInterval.Build(),
		IdleSessionTimeout:       option.IdleSessionTimeout.Build(),
		MinIdleSession:           option.MinIdleSession,
		DialOut:                  anyTlsClient.dialOut,
	})
	if err == nil {
		anyTlsClient.client = client
	}

	return anyTlsClient

}
