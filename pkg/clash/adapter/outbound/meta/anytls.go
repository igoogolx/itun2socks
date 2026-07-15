package clashMeta

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"time"

	clashOutbound "github.com/igoogolx/itun2socks/pkg/clash/adapter/outbound"
	"github.com/igoogolx/itun2socks/pkg/clash/component/dialer"
	clashC "github.com/igoogolx/itun2socks/pkg/clash/constant"
	N "github.com/metacubex/mihomo/common/net"
	C "github.com/metacubex/mihomo/constant"
	"github.com/metacubex/mihomo/transport/anytls"
	"github.com/metacubex/mihomo/transport/vmess"

	M "github.com/metacubex/sing/common/metadata"
	"github.com/metacubex/sing/common/uot"

	metaOutbound "github.com/metacubex/mihomo/adapter/outbound"
)

type AnyTLS struct {
	*metaOutbound.Base
	client             *anytls.Client
	option             *AnyTLSOption
	dialConnContextKey connContextKey
}

type AnyTLSOption struct {
	metaOutbound.BasicOption
	Name                     string                  `proxy:"name"`
	Server                   string                  `proxy:"server"`
	Port                     int                     `proxy:"port"`
	Password                 string                  `proxy:"password"`
	ALPN                     []string                `proxy:"alpn,omitempty"`
	SNI                      string                  `proxy:"sni,omitempty"`
	ECHOpts                  metaOutbound.ECHOptions `proxy:"ech-opts,omitempty"`
	ClientFingerprint        string                  `proxy:"client-fingerprint,omitempty"`
	SkipCertVerify           bool                    `proxy:"skip-cert-verify,omitempty"`
	Fingerprint              string                  `proxy:"fingerprint,omitempty"`
	Certificate              string                  `proxy:"certificate,omitempty"`
	PrivateKey               string                  `proxy:"private-key,omitempty"`
	UDP                      bool                    `proxy:"udp,omitempty"`
	IdleSessionCheckInterval int                     `proxy:"idle-session-check-interval,omitempty"`
	IdleSessionTimeout       int                     `proxy:"idle-session-timeout,omitempty"`
	MinIdleSession           int                     `proxy:"min-idle-session,omitempty"`
}

func (t *AnyTLS) Type() clashC.AdapterType {
	return clashC.AnyTls
}

func (t *AnyTLS) StreamConn(c net.Conn, m *clashC.Metadata) (net.Conn, error) {

	metadata, err := convertMeta(m)
	if err != nil {
		return nil, err

	}

	proxyC, err := t.client.CreateProxy(context.WithValue(context.Background(), t.dialConnContextKey, c), M.ParseSocksaddrHostPort(metadata.String(), metadata.DstPort))
	if err != nil {
		return nil, err
	}

	return proxyC, nil

}

func (t *AnyTLS) Unwrap(_ *clashC.Metadata) clashC.Proxy {
	return nil
}

func (t *AnyTLS) DialContext(ctx context.Context, m *clashC.Metadata, opts ...dialer.Option) (_ clashC.Conn, err error) {

	metadata, err := convertMeta(m)
	if err != nil {
		return nil, err

	}

	dialOutConn, err := dialer.DialContext(ctx, "tcp", t.Addr(), opts...)
	if err != nil {
		return nil, fmt.Errorf("%s connect error: %w", t.Addr(), err)
	}

	c, err := t.client.CreateProxy(context.WithValue(ctx, t.dialConnContextKey, dialOutConn), M.ParseSocksaddrHostPort(metadata.String(), metadata.DstPort))
	if err != nil {
		return nil, err
	}

	return clashOutbound.NewConn(c, t), nil
}

func (t *AnyTLS) ListenPacketContext(ctx context.Context, m *clashC.Metadata, opts ...dialer.Option) (_ clashC.PacketConn, err error) {

	metadata, err := convertMeta(m)
	if err != nil {
		return nil, err

	}

	if err = t.ResolveUDP(ctx, metadata); err != nil {
		return nil, err
	}

	dialOutPc, err := dialer.ListenPacket(ctx, "udp", "", opts...)
	if err != nil {
		return nil, err
	}

	// create tcp
	c, err := t.client.CreateProxy(context.WithValue(ctx, t.dialConnContextKey, dialOutPc), uot.RequestDestination(2))
	if err != nil {
		return nil, err
	}

	// create uot on tcp
	destination := M.SocksaddrFromNet(metadata.UDPAddr())
	return clashOutbound.NewPacketConn(N.NewThreadSafePacketConn(uot.NewLazyConn(c, uot.Request{Destination: destination})), t), nil
}

// SupportUOT implements C.ProxyAdapter
func (t *AnyTLS) SupportUOT() bool {
	return true
}

// ProxyInfo implements C.ProxyAdapter
func (t *AnyTLS) ProxyInfo() C.ProxyInfo {
	info := t.Base.ProxyInfo()
	info.DialerProxy = t.option.DialerProxy
	return info
}

// Close implements C.ProxyAdapter
func (t *AnyTLS) Close() error {
	return t.client.Close()
}

func NewAnyTLS(option AnyTLSOption) (*AnyTLS, error) {
	addr := net.JoinHostPort(option.Server, strconv.Itoa(option.Port))
	outbound := &AnyTLS{
		Base: metaOutbound.NewBase(metaOutbound.BaseOption{
			Name:         option.Name,
			Addr:         addr,
			Type:         C.AnyTLS,
			ProviderName: option.ProviderName,
			UDP:          option.UDP,
			TFO:          option.TFO,
			MPTCP:        option.MPTCP,
			Interface:    option.Interface,
			RoutingMark:  option.RoutingMark,
			Prefer:       option.IPVersion,
		}),
		option:             &option,
		dialConnContextKey: connContextKey("dialOptionsConnKey"),
	}

	tOption := anytls.ClientConfig{
		Password: option.Password,
		Server:   M.ParseSocksaddrHostPort(option.Server, uint16(option.Port)),
		Dialer: AnyTLSDialer{
			addr:          addr,
			diaOutConnKey: outbound.dialConnContextKey,
		},
		IdleSessionCheckInterval: time.Duration(option.IdleSessionCheckInterval) * time.Second,
		IdleSessionTimeout:       time.Duration(option.IdleSessionTimeout) * time.Second,
		MinIdleSession:           option.MinIdleSession,
	}
	echConfig, err := option.ECHOpts.Parse()
	if err != nil {
		return nil, err
	}
	tlsConfig := &vmess.TLSConfig{
		Host:              option.SNI,
		SkipCertVerify:    option.SkipCertVerify,
		NextProtos:        option.ALPN,
		FingerPrint:       option.Fingerprint,
		Certificate:       option.Certificate,
		PrivateKey:        option.PrivateKey,
		ClientFingerprint: option.ClientFingerprint,
		ECH:               echConfig,
	}
	if tlsConfig.Host == "" {
		tlsConfig.Host = option.Server
	}
	tOption.TLSConfig = tlsConfig

	client := anytls.NewClient(context.TODO(), tOption)
	outbound.client = client

	return outbound, nil
}
