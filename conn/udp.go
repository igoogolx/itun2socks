package conn

import (
	"context"
	"github.com/Dreamacro/clash/component/dialer"
	"github.com/Dreamacro/clash/constant"
	C "github.com/Dreamacro/clash/constant"
	"github.com/gofrs/uuid"
	"github.com/igoogolx/itun2socks/constants"
	"github.com/igoogolx/itun2socks/global"
	"net"
	"sync"
	"time"
)

type UdpConn interface {
	ReadFrom([]byte) (int, net.Addr, error)
	WriteTo([]byte, net.Addr) (int, error)
	Close() error
	SetDeadline(t time.Time) error
	SetReadDeadline(t time.Time) error
}

type UdpConnContext struct {
	wg        *sync.WaitGroup
	ctx       context.Context
	id        uuid.UUID
	metadata  *constant.Metadata
	conn      UdpConn
	rule      constants.IpRule
	proxyAddr string
}

func (u *UdpConnContext) ProxyAddr() string {
	return u.proxyAddr
}

func (u *UdpConnContext) Wg() *sync.WaitGroup {
	return u.wg
}

func (u *UdpConnContext) Ctx() context.Context {
	return u.ctx
}

func (u *UdpConnContext) Rule() constants.IpRule {
	return u.rule
}

func (u *UdpConnContext) ID() uuid.UUID {
	return u.id
}

func (u *UdpConnContext) Metadata() *constant.Metadata {
	return u.metadata
}

func (u *UdpConnContext) Conn() UdpConn {
	return u.conn
}

func NewUdpConnContext(ctx context.Context, conn UdpConn, metadata *constant.Metadata, wg *sync.WaitGroup) (*UdpConnContext, error) {
	id, _ := uuid.NewV4()
	rule := global.GetMatcher().GetRule(metadata.DstIP.String())
	proxyAddr, _, err := net.SplitHostPort(getProxy().Addr())
	if err != nil {
		return nil, err
	}
	return &UdpConnContext{
		wg,
		ctx,
		id,
		metadata,
		conn,
		rule,
		proxyAddr,
	}, nil
}

func NewUdpConn(ctx context.Context, metadata *C.Metadata, rule constants.IpRule, defaultInterface string) (net.PacketConn, error) {
	if rule == constants.DistributionBypass {
		return dialer.ListenPacket(ctx, "udp", "", dialer.WithAddrReuse(true), dialer.WithInterface(defaultInterface))
	}
	return getProxy().ListenPacketContext(ctx, metadata, dialer.WithInterface(defaultInterface))
}
