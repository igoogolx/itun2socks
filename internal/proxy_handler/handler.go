package proxy_handler

import (
	"context"
	"github.com/igoogolx/itun2socks/internal/conn"
	"github.com/igoogolx/itun2socks/internal/tunnel"
	"github.com/igoogolx/itun2socks/pkg/log"
	"github.com/sagernet/sing/common/buf"
	"github.com/sagernet/sing/common/bufio"
	"github.com/sagernet/sing/common/bufio/deadline"
	M "github.com/sagernet/sing/common/metadata"
	"github.com/sagernet/sing/common/network"
	"net"
	"sync"
)

type udpConn struct {
	network.PacketConn
}

func (uc *udpConn) ReadFrom(data []byte) (int, net.Addr, error) {

	var err error
	var buff *buf.Buffer
	var dest M.Socksaddr

	defer func() {
		if buff != nil {
			buff.Release()
		}

	}()

	newBuffer := func() *buf.Buffer {
		buff = buf.NewPacket() // do not use stack buffer
		return buff
	}

	dest, err = uc.ReadPacket(newBuffer())

	if err != nil {
		return 0, nil, err
	}

	n, err := buff.Read(data)

	if err != nil {
		return 0, nil, err
	}

	return n, dest, nil
}

func (uc *udpConn) WriteTo(data []byte, addr net.Addr) (int, error) {
	newBuf := buf.NewPacket()
	defer newBuf.Release()
	_, err := newBuf.Write(data)
	if err != nil {
		return 0, err
	}
	err = uc.WritePacket(newBuf, M.SocksaddrFromNet(addr))
	return len(data), err
}

type ConnHandler struct {
	tcpIn chan conn.TcpConnContext
	udpIn chan conn.UdpConnContext
}

func (uc ConnHandler) NewConnection(ctx context.Context, netConn net.Conn, metadata M.Metadata) error {
	local, err := net.ResolveTCPAddr("tcp", metadata.Source.String())
	if err != nil {
		return err
	}
	remote, err := net.ResolveTCPAddr("tcp", metadata.Destination.String())
	if err != nil {
		return err
	}

	m := tunnel.CreateTcpMetadata(*local, *remote)
	var wg sync.WaitGroup
	wg.Add(1)
	ct, err := conn.NewTcpConnContext(ctx, netConn, &m, &wg)
	if err != nil {
		return err
	}
	uc.tcpIn <- *ct
	wg.Wait()
	return nil
}

func (uc ConnHandler) NewPacketConnection(ctx context.Context, packetConn network.PacketConn, metadata M.Metadata) error {
	local, err := net.ResolveUDPAddr("udp", metadata.Source.String())
	if err != nil {
		return err
	}
	remote, err := net.ResolveUDPAddr("udp", metadata.Destination.String())
	if err != nil {
		return err
	}
	m := tunnel.CreateUdpMetadata(*local, *remote)

	if deadline.NeedAdditionalReadDeadline(packetConn) {
		packetConn = deadline.NewFallbackPacketConn(bufio.NewNetPacketConn(packetConn)) // conn from sing should check NeedAdditionalReadDeadline
	}

	var wg sync.WaitGroup
	wg.Add(1)

	ct, err := conn.NewUdpConnContext(ctx, &udpConn{PacketConn: packetConn}, &m, &wg)
	if err != nil {
		return err
	}
	uc.udpIn <- *ct
	wg.Wait()
	return nil
}

func (uc ConnHandler) NewError(_ context.Context, err error) {
	log.Errorln(log.FormatLog(log.TunPrefix, "err: %v"), err)
}
func New(tcpIn chan conn.TcpConnContext,
	udpIn chan conn.UdpConnContext) *ConnHandler {
	return &ConnHandler{
		tcpIn,
		udpIn,
	}
}
