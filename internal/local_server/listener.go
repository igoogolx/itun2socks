package local_server

import (
	"context"
	"github.com/Dreamacro/clash/adapter/inbound"
	C "github.com/Dreamacro/clash/constant"
	"github.com/Dreamacro/clash/listener/mixed"
	"github.com/Dreamacro/clash/listener/socks"
	"github.com/igoogolx/itun2socks/internal/conn"
	"github.com/igoogolx/itun2socks/internal/tunnel"
	"github.com/sagernet/sing/common/buf"
	M "github.com/sagernet/sing/common/metadata"
	"net"
	"sync"
	"time"
)

type udpConn struct {
	C.UDPPacket
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

var tcpIn = make(chan C.ConnContext, 16)
var udpIn = make(chan *inbound.PacketAdapter, 16)

type Listener struct {
	Addr        string
	tcpListener C.Listener
	udpListener C.Listener
}

func process() {
	for t := range tcpIn {
		go processTcp(t)
	}

	for u := range udpIn {
		go processUdp(u)
	}
}

func processTcp(t C.ConnContext) {
	var wg sync.WaitGroup
	wg.Add(1)
	ct, err := conn.NewTcpConnContext(context.Background(), t.Conn(), t.Metadata(), &wg)
	if err != nil {
		return
	}
	tunnel.TcpQueue() <- *ct
	wg.Wait()
}

func processUdp(u *inbound.PacketAdapter) {
	var wg sync.WaitGroup
	wg.Add(1)
	defer wg.Wait()
	ct, err := conn.NewUdpConnContext(context.Background(), udpConn{u.UDPPacket}, u.Metadata(), &wg)
	if err != nil {
		return
	}
	tunnel.UdpQueue() <- *ct
	return
}

func init() {
	go process()
}

func (l *Listener) Start() error {
	tcpListener, err := mixed.New(l.Addr, tcpIn)
	if err != nil {
		return err
	}
	l.tcpListener = tcpListener
	udpListener, err := socks.NewUDP(l.Addr, udpIn)
	if err != nil {
		return err
	}
	l.udpListener = udpListener
	return nil
}

func (l *Listener) Close() error {
	err := l.udpListener.Close()
	if err != nil {
		return err
	}
	err = l.tcpListener.Close()
	if err != nil {
		return err
	}
	return nil
}

func NewListener(httpAddr string) Listener {
	return Listener{Addr: httpAddr}
}
