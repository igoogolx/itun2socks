package proxy_handler

import (
	"context"
	"github.com/igoogolx/itun2socks/internal/conn"
	"github.com/igoogolx/itun2socks/internal/tunnel"
	"github.com/igoogolx/itun2socks/pkg/log"
	sTun "github.com/sagernet/sing-tun"
	"github.com/sagernet/sing/common/bufio"
	"github.com/sagernet/sing/common/bufio/deadline"
	M "github.com/sagernet/sing/common/metadata"
	"github.com/sagernet/sing/common/network"
	"net"
	"net/netip"
	"sync"
)

type ConnHandler struct {
	tcpIn chan conn.TcpConnContext
	udpIn chan conn.UdpConnContext
}

func (uc ConnHandler) JudgeFlow(_ uint8, _ netip.AddrPort, _ netip.AddrPort, _ []byte) sTun.FlowVerdict {
	return sTun.FlowVerdict{
		Action:      0,
		Port:        nil,
		Destination: netip.AddrPort{},
		UDPTimeout:  0,
		NewTracker:  nil,
	}
}

func (uc ConnHandler) NewDNSPacket(_ []byte, _ M.Socksaddr, _ M.Socksaddr, _ network.PacketWriter) {
}

func (uc ConnHandler) NewConnectionEx(ctx context.Context, netConn net.Conn, source M.Socksaddr, destination M.Socksaddr, _ network.CloseHandlerFunc) {
	//TODO: IPV6
	if source.Addr.Is6() {
		return
	}

	m := tunnel.CreateTcpMetadata(source.AddrPort(), destination.AddrPort())
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

	m := tunnel.CreateUdpMetadata(source.AddrPort(), destination.AddrPort())

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
