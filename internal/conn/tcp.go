package conn

import (
	"context"
	"github.com/igoogolx/itun2socks/internal/meta_patch"
	"net"
	"sync"

	"github.com/igoogolx/itun2socks/internal/cfg/distribution/rule_engine"
	C "github.com/igoogolx/itun2socks/pkg/clash/constant"
)

type TcpConnContext struct {
	ctx      context.Context
	metadata *C.Metadata
	conn     net.Conn
	rule     rule_engine.Rule
	wg       *sync.WaitGroup
}

func (t *TcpConnContext) Wg() *sync.WaitGroup {
	return t.wg
}

func (t *TcpConnContext) Ctx() context.Context {
	return t.ctx
}

func (t *TcpConnContext) Rule() rule_engine.Rule {
	return t.rule
}

func (t *TcpConnContext) Metadata() *C.Metadata {
	return t.metadata
}

func (t *TcpConnContext) Conn() net.Conn {
	return t.conn
}

func NewTcpConnContext(ctx context.Context, conn net.Conn, metadata *C.Metadata, wg *sync.WaitGroup) (*TcpConnContext, error) {

	rule := handleMetadata(metadata)

	var connContext = &TcpConnContext{
		ctx,
		metadata,
		conn,
		rule,
		wg,
	}
	return connContext, nil

}

func NewTcpConn(ctx context.Context, metadata *C.Metadata, rule rule_engine.Rule) (net.Conn, error) {
	connDialer, err := GetProxy(rule.GetPolicy())
	if err != nil {
		return nil, err
	}
	return connDialer.DialContext(ctx, meta_patch.ConvertMeta(metadata))
}
