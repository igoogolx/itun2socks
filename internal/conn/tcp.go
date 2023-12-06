package conn

import (
	"context"
	"github.com/Dreamacro/clash/component/dialer"
	C "github.com/Dreamacro/clash/constant"
	"github.com/igoogolx/itun2socks/internal/constants"
	"net"
)

type TcpConnContext struct {
	ctx      context.Context
	metadata *C.Metadata
	conn     net.Conn
	rule     constants.RuleType
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

func NewTcpConnContext(ctx context.Context, conn net.Conn, metadata *C.Metadata) (*TcpConnContext, error) {

	var connContext = &TcpConnContext{
		ctx,
		metadata,
		conn,
		constants.RuleProxy,
	}

	for _, matcher := range GetConnMatcher() {
		rule, err := matcher(metadata, connContext.rule)
		if err == nil {
			connContext.rule = rule
		}
	}

	return connContext, nil

}

func NewTcpConn(ctx context.Context, metadata *C.Metadata, rule constants.RuleType, defaultInterface string) (net.Conn, error) {
	connDialer, err := GetProxy(rule)
	if err != nil {
		return nil, err
	}
	return connDialer.DialContext(ctx, metadata, dialer.WithInterface(defaultInterface))
}
