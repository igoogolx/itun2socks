package conn

import (
	"context"
	"github.com/Dreamacro/clash/component/dialer"
	C "github.com/Dreamacro/clash/constant"
	"github.com/igoogolx/itun2socks/internal/constants"
	"github.com/igoogolx/itun2socks/internal/matcher"
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
	SetWriteDeadline(t time.Time) error
}

type UdpConnContext struct {
	wg       *sync.WaitGroup
	ctx      context.Context
	metadata *C.Metadata
	conn     UdpConn
	rule     constants.RuleType
}

func (u *UdpConnContext) Wg() *sync.WaitGroup {
	return u.wg
}

func (u *UdpConnContext) Ctx() context.Context {
	return u.ctx
}

func (u *UdpConnContext) Rule() constants.RuleType {
	return u.rule
}

func (u *UdpConnContext) Metadata() *C.Metadata {
	return u.metadata
}

func (u *UdpConnContext) Conn() UdpConn {
	return u.conn
}

func NewUdpConnContext(ctx context.Context, conn UdpConn, metadata *C.Metadata, wg *sync.WaitGroup) (*UdpConnContext, error) {
	rule := matcher.GetConnMatcher().GetConnRule(*metadata)

	return &UdpConnContext{
		wg,
		ctx,
		metadata,
		conn,
		rule,
	}, nil
}

func NewUdpConn(ctx context.Context, metadata *C.Metadata, rule constants.RuleType, defaultInterface string) (net.PacketConn, error) {
	return GetProxy(rule).ListenPacketContext(ctx, metadata, dialer.WithInterface(defaultInterface), dialer.WithAddrReuse(true))
}
