package conn

import (
	"context"
	"github.com/Dreamacro/clash/component/dialer"
	C "github.com/Dreamacro/clash/constant"
	"github.com/igoogolx/itun2socks/internal/cfg/distribution/ruleEngine"
	"github.com/sagernet/sing/common/buf"
	M "github.com/sagernet/sing/common/metadata"
	"github.com/sagernet/sing/common/network"
	"net"
	"sync"
	"time"
)

type UdpConn interface {
	network.PacketReader
	network.PacketWriter
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

type CopyableUdpConn struct {
	net.PacketConn
}

func (c *CopyableUdpConn) ReadPacket(buffer *buf.Buffer) (destination M.Socksaddr, err error) {
	var rawData []byte
	n, addr, err := c.ReadFrom(rawData)
	if err != nil {
		return M.Socksaddr{}, err
	}
	_, err = buffer.Write(rawData[0:n])
	return M.SocksaddrFromNet(addr), err
}

func (c *CopyableUdpConn) WritePacket(buffer *buf.Buffer, destination M.Socksaddr) error {
	_, err := c.WriteTo(buffer.Bytes(), destination)
	return err
}

func NewUdpConn(ctx context.Context, metadata *C.Metadata, rule ruleEngine.Rule, defaultInterface string) (*CopyableUdpConn, error) {
	connDialer, err := GetProxy(rule.GetPolicy())
	if err != nil {
		return nil, err
	}
	rawConn, err := connDialer.ListenPacketContext(ctx, metadata, dialer.WithInterface(defaultInterface), dialer.WithAddrReuse(true))
	if err != nil {
		return nil, err
	}
	return &CopyableUdpConn{rawConn}, nil
}
