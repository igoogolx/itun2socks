package conn

import (
	"context"
	"fmt"
	"github.com/Dreamacro/clash/component/dialer"
	C "github.com/Dreamacro/clash/constant"
	"github.com/igoogolx/itun2socks/internal/cfg/distribution/rule_engine"
	"github.com/igoogolx/itun2socks/internal/constants"
	"github.com/igoogolx/itun2socks/internal/dns"
	"github.com/igoogolx/itun2socks/pkg/log"
	"github.com/igoogolx/itun2socks/pkg/pool"
	"github.com/sagernet/sing/common/buf"
	E "github.com/sagernet/sing/common/exceptions"
	M "github.com/sagernet/sing/common/metadata"
	"github.com/sagernet/sing/common/network"
	"net"
	"sync"
	"time"
)

type UdpConnContext struct {
	ctx      context.Context
	metadata *C.Metadata
	conn     network.PacketConn
	rule     rule_engine.Rule
	wg       *sync.WaitGroup
}

func (u *UdpConnContext) Wg() *sync.WaitGroup {
	return u.wg
}

func (u *UdpConnContext) Ctx() context.Context {
	return u.ctx
}

func (u *UdpConnContext) Rule() rule_engine.Rule {
	return u.rule
}

func (u *UdpConnContext) Metadata() *C.Metadata {
	return u.metadata
}

func (u *UdpConnContext) Conn() network.PacketConn {
	return u.conn
}

func NewUdpConnContext(ctx context.Context, conn network.PacketConn, metadata *C.Metadata, wg *sync.WaitGroup) (*UdpConnContext, error) {
	rule := resolveMetadata(metadata)

	if rule.GetPolicy() == constants.PolicyProxy && defaultIsFakeIpEnabled {
		hostByFakeIp, ok := dns.FakeIpPool.LookBack(metadata.DstIP)
		if ok {
			metadata.Host = hostByFakeIp
		}
	}
	var connContext = &UdpConnContext{
		ctx,
		metadata,
		conn,
		rule,
		wg,
	}

	return connContext, nil
}

type CopyablePacketConn struct {
	net.PacketConn
}

func shouldIgnorePacketError(err error) bool {
	// ignore simple error
	if E.IsTimeout(err) || E.IsClosed(err) || E.IsCanceled(err) {
		return true
	}
	return false
}

func PrintPacketError(err error, msg string) {
	printLog := log.Warnln
	if shouldIgnorePacketError(err) {
		printLog = log.Debugln
	}
	printLog(msg)
}

func (c *CopyablePacketConn) ReadPacket(buffer *buf.Buffer) (destination M.Socksaddr, err error) {
	receivedBuf := pool.NewBytes(pool.BufSize)
	defer pool.FreeBytes(receivedBuf)
	for {
		err := c.SetReadDeadline(time.Now().Add(5 * time.Second))
		if err != nil {
			return M.Socksaddr{}, fmt.Errorf("fail to set udp conn read deadline: %v", err)
		}
		n, addr, err := c.ReadFrom(receivedBuf)
		if shouldIgnorePacketError(err) {
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

func (c *CopyablePacketConn) WritePacket(buffer *buf.Buffer, destination M.Socksaddr) error {
	_, err := c.WriteTo(buffer.Bytes(), destination)
	return err
}

type CopyableReaderWriterConn struct {
	network.PacketConn
}

func (uc *CopyableReaderWriterConn) ReadFrom(data []byte) (int, net.Addr, error) {

	var err error
	var dest M.Socksaddr

	buff := buf.NewPacket()

	defer buff.Release()
	dest, err = uc.ReadPacket(buff)

	if err != nil {
		return 0, nil, err
	}

	n, err := buff.Read(data)

	if err != nil {
		return 0, nil, err
	}

	return n, dest, nil
}

func (uc *CopyableReaderWriterConn) WriteTo(data []byte, addr net.Addr) (int, error) {
	newBuf := buf.NewPacket()
	defer newBuf.Release()
	_, err := newBuf.Write(data)
	if err != nil {
		return 0, err
	}
	err = uc.WritePacket(newBuf, M.SocksaddrFromNet(addr))
	return len(data), err
}

func NewUdpConn(ctx context.Context, metadata *C.Metadata, rule rule_engine.Rule, defaultInterface string) (*CopyablePacketConn, error) {
	connDialer, err := GetProxy(rule.GetPolicy())
	if err != nil {
		return nil, err
	}
	rawConn, err := connDialer.ListenPacketContext(ctx, metadata, dialer.WithInterface(defaultInterface), dialer.WithAddrReuse(true))
	if err != nil {
		return nil, err
	}
	return &CopyablePacketConn{rawConn}, nil
}
