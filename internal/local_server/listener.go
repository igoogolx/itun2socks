package local_server

import (
	"context"
	"github.com/igoogolx/itun2socks/internal/meta_patch"
	"github.com/metacubex/mihomo/adapter/inbound"
	"github.com/metacubex/mihomo/listener/mixed"
	"github.com/metacubex/mihomo/listener/socks"
	"net"
	"net/netip"
	"sync"
	"time"

	"github.com/igoogolx/itun2socks/internal/conn"
	"github.com/igoogolx/itun2socks/internal/tunnel"
	metaC "github.com/metacubex/mihomo/constant"
	"github.com/sagernet/sing/common/buf"
	M "github.com/sagernet/sing/common/metadata"
)

type udpConn struct {
	metaC.UDPPacket
}

func (u udpConn) ReadPacket(buffer *buf.Buffer) (destination M.Socksaddr, err error) {
	_, addr, err := u.readFrom(buffer.Bytes())
	return M.SocksaddrFromNet(addr), err
}

func (u udpConn) WritePacket(buffer *buf.Buffer, destination M.Socksaddr) error {
	_, err := u.writeTo(buffer.Bytes(), destination)
	return err
}

func (u udpConn) readFrom(data []byte) (int, net.Addr, error) {
	n := copy(data, u.Data())
	return n, u.LocalAddr(), nil
}

func (u udpConn) writeTo(data []byte, addr net.Addr) (int, error) {
	n, err := u.WriteBack(data, addr)
	return n, err
}

func (u udpConn) Close() error {
	u.Drop()
	return nil
}

func (u udpConn) SetDeadline(_ time.Time) error {
	return nil
}

func (u udpConn) SetReadDeadline(_ time.Time) error {
	return nil
}

func (u udpConn) SetWriteDeadline(_ time.Time) error {
	return nil
}

type Listener struct {
	Addr        string
	tcpListener metaC.Listener
	udpListener metaC.Listener
	Port        int
}

type ListenerHandler struct {
}

func (l ListenerHandler) HandleTCPConn(c net.Conn, metadata *metaC.Metadata) {
	var wg sync.WaitGroup
	wg.Add(1)
	ct, err := conn.NewTcpConnContext(context.Background(), c, meta_patch.ConvertClash(metadata), &wg)
	if err != nil {
		return
	}
	tunnel.TcpQueue() <- *ct
	wg.Wait()
}

func (l ListenerHandler) HandleUDPPacket(packet metaC.UDPPacket, metadata *metaC.Metadata) {

	var wg sync.WaitGroup
	wg.Add(1)
	defer wg.Wait()
	ct, err := conn.NewUdpConnContext(context.Background(), udpConn{packet}, meta_patch.ConvertClash(metadata), &wg)
	if err != nil {
		return
	}
	tunnel.UdpQueue() <- *ct
	return

}

func (l ListenerHandler) NatTable() metaC.NatTable {
	return nil
}

var listenerHandler = ListenerHandler{}

func (l *Listener) Start() error {
	inbound.SetAllowedIPs([]netip.Prefix{
		netip.MustParsePrefix("0.0.0.0/0"),
	})
	tcpListener, err := mixed.New(l.Addr, listenerHandler)
	if err != nil {
		return err
	}
	l.tcpListener = tcpListener
	udpListener, err := socks.NewUDP(l.Addr, listenerHandler)
	if err != nil {
		return err
	}
	l.udpListener = udpListener
	return nil
}

func (l *Listener) Close() error {
	if l.udpListener != nil {
		err := l.udpListener.Close()
		if err != nil {
			return err
		}
	}

	if l.tcpListener != nil {
		err := l.tcpListener.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

func NewListener(httpAddr string, port int) Listener {
	return Listener{Addr: httpAddr, Port: port}
}
