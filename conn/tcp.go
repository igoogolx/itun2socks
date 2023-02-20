package conn

import (
	"context"
	"github.com/Dreamacro/clash/component/dialer"
	"github.com/Dreamacro/clash/constant"
	C "github.com/Dreamacro/clash/constant"
	"github.com/gofrs/uuid"
	"github.com/igoogolx/itun2socks/constants"
	"net"
	"sync"
)

type TcpConnContext struct {
	wg       *sync.WaitGroup
	ctx      context.Context
	id       uuid.UUID
	metadata *constant.Metadata
	conn     net.Conn
	rule     constants.IpRule
}

func (t *TcpConnContext) Wg() *sync.WaitGroup {
	return t.wg
}

func (t *TcpConnContext) Ctx() context.Context {
	return t.ctx
}

func (t *TcpConnContext) Rule() constants.IpRule {
	return t.rule
}

func (t *TcpConnContext) ID() uuid.UUID {
	return t.id
}

func (t *TcpConnContext) Metadata() *constant.Metadata {
	return t.metadata
}

func (t *TcpConnContext) Conn() net.Conn {
	return t.conn
}

func NewTcpConnContext(ctx context.Context, conn net.Conn, metadata *constant.Metadata, wg *sync.WaitGroup) *TcpConnContext {
	id, _ := uuid.NewV4()
	rule := GetMatcher().GetRule(metadata.DstIP.String())
	return &TcpConnContext{
		wg,
		ctx,
		id,
		metadata,
		conn,
		rule,
	}

}

func NewTcpConn(ctx context.Context, metadata *C.Metadata, rule constants.IpRule, defaultInterface string) (net.Conn, error) {
	if rule == constants.DistributionBypass {
		return dialer.DialContext(ctx, "tcp", metadata.RemoteAddress(), dialer.WithInterface(defaultInterface))
	}
	return getProxy().DialContext(ctx, metadata, dialer.WithInterface(defaultInterface))
}
