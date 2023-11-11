package conn

import (
	"context"
	"github.com/Dreamacro/clash/component/dialer"
	C "github.com/Dreamacro/clash/constant"
	"github.com/igoogolx/itun2socks/internal/constants"
	"github.com/igoogolx/itun2socks/internal/matcher"
	"net"
	"sync"
)

type TcpConnContext struct {
	wg       *sync.WaitGroup
	ctx      context.Context
	metadata *C.Metadata
	conn     net.Conn
	rule     constants.RuleType
}

func (t *TcpConnContext) Wg() *sync.WaitGroup {
	return t.wg
}

func (t *TcpConnContext) Ctx() context.Context {
	return t.ctx
}

func (t *TcpConnContext) Rule() constants.RuleType {
	return t.rule
}

func (t *TcpConnContext) Metadata() *C.Metadata {
	return t.metadata
}

func (t *TcpConnContext) Conn() net.Conn {
	return t.conn
}

func NewTcpConnContext(ctx context.Context, conn net.Conn, metadata *C.Metadata, wg *sync.WaitGroup) (*TcpConnContext, error) {
	rule := matcher.GetConnMatcher().GetConnRule(*metadata)
	return &TcpConnContext{
		wg,
		ctx,
		metadata,
		conn,
		rule,
	}, nil

}

func NewTcpConn(ctx context.Context, metadata *C.Metadata, rule constants.RuleType, defaultInterface string) (net.Conn, error) {
	return GetProxy(rule).DialContext(ctx, metadata, dialer.WithInterface(defaultInterface))
}
