package clashMeta

import (
	"context"
	"fmt"
	"net"
	"net/netip"

	clashC "github.com/igoogolx/itun2socks/pkg/clash/constant"
	metaC "github.com/metacubex/mihomo/constant"

	M "github.com/metacubex/sing/common/metadata"
)

type connContextKey string

type AnyTLSDialer struct {
	diaOutConnKey connContextKey
	addr          string
}

func (a AnyTLSDialer) DialContext(ctx context.Context, _ string, _ M.Socksaddr) (net.Conn, error) {

	v, ok := ctx.Value(a.diaOutConnKey).(net.Conn)

	if !ok {
		return nil, fmt.Errorf("invalid dialer")
	}

	return v, nil

}

func (a AnyTLSDialer) ListenPacket(ctx context.Context, _ M.Socksaddr) (net.PacketConn, error) {

	v, ok := ctx.Value(a.diaOutConnKey).(net.PacketConn)

	if !ok {
		return nil, fmt.Errorf("invalid dialer")
	}

	return v, nil

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
