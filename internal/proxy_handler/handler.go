package proxy_handler

import (
	"context"
	"github.com/igoogolx/itun2socks/internal/conn"
	"github.com/igoogolx/itun2socks/internal/tunnel"
	"github.com/igoogolx/itun2socks/pkg/log"
	"github.com/sagernet/sing/common/bufio"
	"github.com/sagernet/sing/common/bufio/deadline"
	M "github.com/sagernet/sing/common/metadata"
	"github.com/sagernet/sing/common/network"
	"net"
	"sync"
)

type ConnHandler struct {
	tcpIn chan conn.TcpConnContext
	udpIn chan conn.UdpConnContext
}

func (uc ConnHandler) PrepareConnection(_ string, _ M.Socksaddr, _ M.Socksaddr) error {
	return nil
}

func (uc ConnHandler) NewConnectionEx(ctx context.Context, netConn net.Conn, source M.Socksaddr, destination M.Socksaddr, _ network.CloseHandlerFunc) {
	//TODO: IPV6
	if source.Addr.Is6() {
		return
	}
	local, err := net.ResolveTCPAddr("tcp", source.String())
	if err != nil {
		return
	}
	remote, err := net.ResolveTCPAddr("tcp", destination.String())
	if err != nil {
		return
	}

	m := tunnel.CreateTcpMetadata(*local, *remote)
	var wg sync.WaitGroup
	wg.Add(1)
	ct, err := conn.NewTcpConnContext(ctx, netConn, &m, &wg)
	if err != nil {
		return
	}
	uc.tcpIn <- *ct
	wg.Wait()
}

func (uc ConnHandler) NewPacketConnectionEx(ctx context.Context, packetConn network.PacketConn, source M.Socksaddr, destination M.Socksaddr, _ network.CloseHandlerFunc) {
	//TODO: IPV6
	if source.Addr.Is6() {
		return
	}
	local, err := net.ResolveUDPAddr("udp", source.String())
	if err != nil {
		return
	}
	remote, err := net.ResolveUDPAddr("udp", destination.String())
	if err != nil {
		return
	}
	m := tunnel.CreateUdpMetadata(*local, *remote)

	if deadline.NeedAdditionalReadDeadline(packetConn) {
		packetConn = deadline.NewFallbackPacketConn(bufio.NewNetPacketConn(packetConn)) // conn from sing should check NeedAdditionalReadDeadline
	}

	var wg sync.WaitGroup
	wg.Add(1)

	ct, err := conn.NewUdpConnContext(ctx, packetConn, &m, &wg)
	if err != nil {
		return
	}
	uc.udpIn <- *ct
	wg.Wait()
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
