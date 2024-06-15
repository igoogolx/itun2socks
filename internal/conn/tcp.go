package conn

import (
	"context"
	"github.com/Dreamacro/clash/component/dialer"
	C "github.com/Dreamacro/clash/constant"
	"github.com/igoogolx/itun2socks/internal/cfg/distribution/ruleEngine"
	"net"
	"sync"
)

type TcpConnContext struct {
	ctx      context.Context
	metadata *C.Metadata
	conn     net.Conn
	rule     ruleEngine.Rule
	wg       *sync.WaitGroup
}

func (t *TcpConnContext) Wg() *sync.WaitGroup {
	return t.wg
}

func (t *TcpConnContext) Ctx() context.Context {
	return t.ctx
}

func (t *TcpConnContext) Rule() ruleEngine.Rule {
	return t.rule
}

func (t *TcpConnContext) Metadata() *C.Metadata {
	return t.metadata
}

func (t *TcpConnContext) Conn() net.Conn {
	return t.conn
}

func NewTcpConnContext(ctx context.Context, conn net.Conn, metadata *C.Metadata, wg *sync.WaitGroup) (*TcpConnContext, error) {

	var connContext = &TcpConnContext{
		ctx,
		metadata,
		conn,
		ruleEngine.BuiltInProxyRule,
		wg,
	}

	for _, matcher := range GetConnMatcher() {
		rule, err := matcher(metadata, connContext.rule)
		if err == nil {
			connContext.rule = rule
		}
	}

	return connContext, nil

}

func NewTcpConn(ctx context.Context, metadata *C.Metadata, rule ruleEngine.Rule, defaultInterface string) (net.Conn, error) {
	connDialer, err := GetProxy(rule.GetPolicy())
	if err != nil {
		return nil, err
	}
	return connDialer.DialContext(ctx, metadata, dialer.WithInterface(defaultInterface))
}
