package proxy_handler

import (
	"context"
	"github.com/Dreamacro/clash/log"
	conn2 "github.com/igoogolx/itun2socks/conn"
	tunnel2 "github.com/igoogolx/itun2socks/tunnel"
	buf2 "github.com/sagernet/sing/common/buf"
	M "github.com/sagernet/sing/common/metadata"
	"github.com/sagernet/sing/common/network"
	"net"
	"sync"
)

type udpConn struct {
	network.PacketConn
}

func (uc udpConn) ReadFrom(buf []byte) (int, net.Addr, error) {
	newBuf := buf2.NewPacket()
	defer newBuf.Release()
	des, err := uc.ReadPacket(newBuf)
	if err != nil {
		return 0, nil, err
	}
	n, err := newBuf.Read(buf)
	if err != nil {
		return 0, nil, err
	}
	return n, des.UDPAddr(), nil
}

func (uc udpConn) WriteTo(buf []byte, addr net.Addr) (int, error) {
	newBuf := buf2.NewPacket()
	_, err := newBuf.Write(buf)
	defer newBuf.Release()
	if err != nil {
		return 0, err
	}
	err = uc.WritePacket(newBuf, M.SocksaddrFromNet(addr))
	return len(buf), err
}

type ConnHandler struct {
}

func (c ConnHandler) NewConnection(ctx context.Context, conn net.Conn, metadata M.Metadata) error {
	local, err := net.ResolveTCPAddr("tcp", metadata.Source.String())
	if err != nil {
		return err
	}
	remote, err := net.ResolveTCPAddr("tcp", metadata.Destination.String())
	if err != nil {
		return err
	}
	m := tunnel2.CreateTcpMetadata(*local, *remote)
	var wg sync.WaitGroup
	wg.Add(1)
	defer wg.Wait()
	ct := conn2.NewTcpConnContext(ctx, conn, &m, &wg)
	tunnel2.AddTcpConn(*ct)
	return nil
}

func (c ConnHandler) NewPacketConnection(ctx context.Context, conn network.PacketConn, metadata M.Metadata) error {
	local, err := net.ResolveUDPAddr("udp", metadata.Source.String())
	if err != nil {
		return err
	}
	remote, err := net.ResolveUDPAddr("udp", metadata.Destination.String())
	if err != nil {
		return err
	}
	m := tunnel2.CreateUdpMetadata(*local, *remote)
	var wg sync.WaitGroup
	wg.Add(1)
	defer wg.Wait()
	ct, err := conn2.NewUdpConnContext(ctx, udpConn{conn}, &m, &wg)
	if err != nil {
		return err
	}
	tunnel2.AddUdpConn(ct)
	return nil
}

func (c ConnHandler) NewError(ctx context.Context, err error) {
	log.Errorln("proxy handler, err: %v", err)
}
func New() *ConnHandler {
	return &ConnHandler{}
}
