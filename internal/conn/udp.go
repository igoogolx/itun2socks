package conn

import (
	"context"
	"github.com/Dreamacro/clash/component/dialer"
	C "github.com/Dreamacro/clash/constant"
	"github.com/igoogolx/itun2socks/internal/cfg/distribution/ruleEngine"
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
	ctx      context.Context
	metadata *C.Metadata
	conn     UdpConn
	rule     ruleEngine.Rule
	wg       *sync.WaitGroup
}

func (u *UdpConnContext) Wg() *sync.WaitGroup {
	return u.wg
}

func (u *UdpConnContext) Ctx() context.Context {
	return u.ctx
}

func (u *UdpConnContext) Rule() ruleEngine.Rule {
	return u.rule
}

func (u *UdpConnContext) Metadata() *C.Metadata {
	return u.metadata
}

func (u *UdpConnContext) Conn() UdpConn {
	return u.conn
}

func NewUdpConnContext(ctx context.Context, conn UdpConn, metadata *C.Metadata, wg *sync.WaitGroup) (*UdpConnContext, error) {
	var connContext = &UdpConnContext{
		ctx,
		metadata,
		conn,
		resolveMetadata(metadata),
		wg,
	}

	return connContext, nil
}

func NewUdpConn(ctx context.Context, metadata *C.Metadata, rule ruleEngine.Rule, defaultInterface string) (net.PacketConn, error) {
	connDialer, err := GetProxy(rule.GetPolicy())
	if err != nil {
		return nil, err
	}
	return connDialer.ListenPacketContext(ctx, metadata, dialer.WithInterface(defaultInterface), dialer.WithAddrReuse(true))
}
