package conn

import (
	"context"
	"github.com/igoogolx/itun2socks/internal/constants"
	metaC "github.com/metacubex/mihomo/constant"
)

type proxyAdapter struct {
	policy constants.Policy
}

func (a proxyAdapter) Name() string {
	p, _ := GetProxy(a.policy)
	return p.Name()
}

func (a proxyAdapter) Type() metaC.AdapterType {

	p, _ := GetProxy(a.policy)
	return p.Type()
}

func (a proxyAdapter) Addr() string {

	p, _ := GetProxy(a.policy)
	return p.Addr()

}

func (a proxyAdapter) SupportUDP() bool {

	p, _ := GetProxy(a.policy)
	return p.SupportUDP()
}

func (a proxyAdapter) ProxyInfo() metaC.ProxyInfo {

	p, _ := GetProxy(a.policy)
	return p.ProxyInfo()
}

func (a proxyAdapter) MarshalJSON() ([]byte, error) {

	p, _ := GetProxy(a.policy)
	return p.MarshalJSON()
}

func (a proxyAdapter) DialContext(ctx context.Context, metadata *metaC.Metadata) (metaC.Conn, error) {

	p, _ := GetProxy(a.policy)
	return p.DialContext(ctx, metadata)
}

func (a proxyAdapter) ListenPacketContext(ctx context.Context, metadata *metaC.Metadata) (metaC.PacketConn, error) {

	p, _ := GetProxy(a.policy)
	return p.ListenPacketContext(ctx, metadata)
}

func (a proxyAdapter) SupportUOT() bool {
	p, _ := GetProxy(a.policy)
	return p.SupportUOT()
}

func (a proxyAdapter) IsL3Protocol(metadata *metaC.Metadata) bool {
	p, _ := GetProxy(a.policy)
	return p.IsL3Protocol(metadata)
}

func (a proxyAdapter) Unwrap(metadata *metaC.Metadata, touch bool) metaC.Proxy {
	p, _ := GetProxy(a.policy)
	return p.Unwrap(metadata, touch)
}

func (a proxyAdapter) Close() error {
	p, _ := GetProxy(a.policy)
	return p.Close()
}

var ProxyAdapter = proxyAdapter{constants.PolicyProxy}
var DirectAdapter = proxyAdapter{constants.PolicyDirect}
