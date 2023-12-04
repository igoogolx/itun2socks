package proxy_handler

import (
	"context"
	"github.com/igoogolx/itun2socks/internal/conn"
	"github.com/igoogolx/itun2socks/internal/tunnel"
	"github.com/igoogolx/itun2socks/pkg/log"
	"github.com/sagernet/sing/common/buf"
	"github.com/sagernet/sing/common/bufio"
	M "github.com/sagernet/sing/common/metadata"
	"github.com/sagernet/sing/common/network"
	"io"
	"net"
	"time"
)

type udpConn struct {
	network.PacketConn
	read bool
	dest M.Socksaddr
	buff *buf.Buffer
}

func (uc *udpConn) ReadFrom(data []byte) (int, net.Addr, error) {
	if uc.read {
		return 0, nil, io.EOF
	}
	n, err := uc.buff.Read(data)
	uc.buff.Release()
	uc.read = true
	return n, uc.dest.UDPAddr(), err
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

func (uc ConnHandler) SetReadDeadline(t time.Time) error {
	return nil
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
	ct, err := conn.NewTcpConnContext(ctx, netConn, &m)
	if err != nil {
		return err
	}
	uc.tcpIn <- *ct
	return nil
}

func (uc ConnHandler) NewPacketConnection(ctx context.Context, packetConn network.PacketConn, metadata M.Metadata) error {
	defer func(packetConn network.PacketConn) {
		err := packetConn.Close()
		if err != nil {
			log.Errorln("fail to close packetConn")
		}
	}(packetConn)
	local, err := net.ResolveUDPAddr("udp", metadata.Source.String())
	if err != nil {
		return err
	}
	remote, err := net.ResolveUDPAddr("udp", metadata.Destination.String())
	if err != nil {
		return err
	}
	m := tunnel.CreateUdpMetadata(*local, *remote)

	for {
		var buff *buf.Buffer
		newBuffer := func() *buf.Buffer {
			buff = buf.NewPacket() // do not use stack buffer
			return buff
		}
		var err error
		var dest M.Socksaddr
		readWaiter, isReadWaiter := bufio.CreatePacketReadWaiter(packetConn)
		if isReadWaiter {
			readWaiter.InitializeReadWaiter(newBuffer)
		}
		if isReadWaiter {
			dest, err = readWaiter.WaitReadPacket()
		} else {
			dest, err = packetConn.ReadPacket(newBuffer())
		}
		if err != nil {
			if buff != nil {
				buff.Release()
			}
			break
		}
		ct, err := conn.NewUdpConnContext(ctx, &udpConn{PacketConn: packetConn, dest: dest, buff: buff}, &m)
		if err != nil {
			if buff != nil {
				buff.Release()
			}
			break
		}
		uc.udpIn <- *ct
	}

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
