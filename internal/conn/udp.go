package conn

import (
	"context"
	"github.com/Dreamacro/clash/component/dialer"
	"github.com/Dreamacro/clash/constant"
	C "github.com/Dreamacro/clash/constant"
	"github.com/gofrs/uuid"
	"github.com/igoogolx/itun2socks/internal/constants"
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
	wg       *sync.WaitGroup
	ctx      context.Context
	id       uuid.UUID
	metadata *constant.Metadata
	conn     UdpConn
	rule     constants.IpRule
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
	rule := GetMatcher().GetRule(metadata.DstIP.String())

	return &UdpConnContext{
		wg,
		ctx,
		id,
		metadata,
		conn,
		rule,
	}, nil
}

func NewUdpConn(ctx context.Context, metadata *C.Metadata, rule constants.IpRule, defaultInterface string) (net.PacketConn, error) {
	return getProxy(rule).ListenPacketContext(ctx, metadata, dialer.WithInterface(defaultInterface))
}
