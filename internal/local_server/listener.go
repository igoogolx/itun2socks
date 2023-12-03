package local_server

import (
	"context"
	"github.com/Dreamacro/clash/adapter/inbound"
	C "github.com/Dreamacro/clash/constant"
	"github.com/Dreamacro/clash/listener/mixed"
	"github.com/Dreamacro/clash/listener/socks"
	"github.com/igoogolx/itun2socks/internal/conn"
	"github.com/igoogolx/itun2socks/internal/tunnel"
	"net"
	"sync"
	"time"
)

type udpConn struct {
	C.UDPPacket
}

func (uc udpConn) ReadFrom(data []byte) (int, net.Addr, error) {
	n := copy(data, uc.Data())
	return n, uc.LocalAddr(), nil
}

func (uc udpConn) WriteTo(data []byte, addr net.Addr) (int, error) {
	n, err := uc.WriteBack(data, addr)
	return n, err
}

func (uc udpConn) Close() error {
	uc.Drop()
	return nil
}

func (uc udpConn) SetDeadline(t time.Time) error {
	return nil
}

func (uc udpConn) SetReadDeadline(t time.Time) error {
	return nil
}

func (uc udpConn) SetWriteDeadline(t time.Time) error {
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
	defer wg.Wait()
	ct, err := conn.NewTcpConnContext(context.Background(), t.Conn(), t.Metadata(), &wg)
	if err != nil {
		return
	}
	tunnel.TcpQueue() <- *ct
}

func processUdp(u *inbound.PacketAdapter) {
	ct, err := conn.NewUdpConnContext(context.Background(), udpConn{u.UDPPacket}, u.Metadata())
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
