package conn

import (
	"context"
	"fmt"
	"github.com/Dreamacro/clash/component/dialer"
	C "github.com/Dreamacro/clash/constant"
	"github.com/igoogolx/itun2socks/internal/cfg/distribution/ruleEngine"
	"github.com/igoogolx/itun2socks/pkg/pool"
	"github.com/sagernet/sing/common/buf"
	E "github.com/sagernet/sing/common/exceptions"
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

func ShouldIgnorePacketError(err error) bool {
	// ignore simple error
	if E.IsTimeout(err) || E.IsClosed(err) || E.IsCanceled(err) {
		return true
	}
	return false
}

func (c *CopyableUdpConn) ReadPacket(buffer *buf.Buffer) (destination M.Socksaddr, err error) {
	receivedBuf := pool.NewBytes(pool.BufSize)
	defer pool.FreeBytes(receivedBuf)
	for {
		err := c.SetReadDeadline(time.Now().Add(5 * time.Second))
		if err != nil {
			return M.Socksaddr{}, fmt.Errorf("fail to set udp conn read deadline: %v", err)
		}
		n, addr, err := c.ReadFrom(receivedBuf)
		if ShouldIgnorePacketError(err) {
			return M.SocksaddrFromNet(addr), nil
		}
		if err != nil {
			return M.Socksaddr{}, fmt.Errorf("fail to read udp from copyable conn:%v", err)
		}
		_, err = buffer.Write(receivedBuf[:n])
		if err != nil {
			return M.Socksaddr{}, fmt.Errorf("fail to write udp to bufffer:%v", err)
		}
	}
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
