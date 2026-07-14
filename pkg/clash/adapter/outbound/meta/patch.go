package clashMeta

import (
	"context"
	"fmt"
	"net"
	"net/netip"

	"github.com/igoogolx/itun2socks/pkg/clash/component/dialer"
	clashC "github.com/igoogolx/itun2socks/pkg/clash/constant"
	metaC "github.com/metacubex/mihomo/constant"

	M "github.com/metacubex/sing/common/metadata"
)

type connOptionsContextKey string

type AnyTLSDialer struct {
	diaOutConnKey connOptionsContextKey
	addr          string
}

func (a AnyTLSDialer) DialContext(ctx context.Context, network string, destination M.Socksaddr) (net.Conn, error) {

	v, ok := ctx.Value(a.diaOutConnKey).([]dialer.Option)

	if !ok {
		return nil, fmt.Errorf("invalid dialer options")
	}

	c, err := dialer.DialContext(ctx, network, a.addr, v...)
	if err != nil {
		return nil, fmt.Errorf("%s connect error: %w", a.addr, err)
	}

	return c, nil

}

func (a AnyTLSDialer) ListenPacket(ctx context.Context, destination M.Socksaddr) (net.PacketConn, error) {

	v, ok := ctx.Value(a.diaOutConnKey).([]dialer.Option)

	if !ok {
		return nil, fmt.Errorf("invalid dialer options")
	}

	pc, err := dialer.ListenPacket(ctx, "udp", "", v...)
	if err != nil {
		return nil, err
	}

	return pc, nil

}

func convertMeta(m *clashC.Metadata) (*metaC.Metadata, error) {

	srcIP, err := netip.ParseAddr(m.SrcIP.String())
	if err != nil {
		srcIP = netip.Addr{}
	}

	dstIP, err := netip.ParseAddr(m.DstIP.String())
	if err != nil {
		dstIP = netip.Addr{}
	}

	newMetadata := &metaC.Metadata{
		NetWork:     metaC.NetWork(m.NetWork),
		Type:        metaC.Type(m.Type),
		SrcIP:       srcIP,
		DstIP:       dstIP,
		SrcPort:     uint16(m.SrcPort),
		DstPort:     uint16(m.DstPort),
		Host:        m.Host,
		ProcessPath: m.ProcessPath,
	}

	return newMetadata, nil

}
