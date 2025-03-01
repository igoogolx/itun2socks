package tunnel

import (
	"fmt"
	"github.com/igoogolx/itun2socks/internal/conn"
	"github.com/igoogolx/itun2socks/internal/dns"
	"github.com/igoogolx/itun2socks/internal/tunnel/statistic"
	"github.com/igoogolx/itun2socks/pkg/log"
	"github.com/igoogolx/itun2socks/pkg/network_iface"
	"github.com/igoogolx/itun2socks/pkg/pool"
	D "github.com/miekg/dns"
	"github.com/sagernet/sing/common/bufio"
	"github.com/sagernet/sing/common/network"
	"sync"
)

var (
	udpQueue = make(chan conn.UdpConnContext, 1024)
)

func UdpQueue() chan conn.UdpConnContext {
	return udpQueue
}

func copyUdpPacket(lc network.PacketConn, rc network.PacketConn) error {
	_, err := bufio.CopyPacket(lc, rc)
	return err
}

func handleUdpConn(ct conn.UdpConnContext) {
	var once sync.Once
	var lc network.PacketConn
	var err error

	cleanConn := func() {
		if lc != nil {
			err := closeConn(lc)
			if err != nil {
				log.Debugln(log.FormatLog(log.UdpPrefix, "fail to close local conn,err: %v"), err)
			}
		}
		err := closeConn(ct.Conn())
		if err != nil {
			log.Warnln(log.FormatLog(log.UdpPrefix, "fail to close remote conn,err: %v"), err)
		}
	}

	defer func() {
		ct.Wg().Done()
		once.Do(cleanConn)
	}()

	localConn, err := conn.NewUdpConn(ct.Ctx(), ct.Metadata(), ct.Rule(), network_iface.GetDefaultInterfaceName())
	if err != nil {
		log.Warnln(log.FormatLog(log.UdpPrefix, "fail to get udp conn, err: %v, remote address: %v"), err, ct.Metadata().RemoteAddress())
		return
	}
	lc = statistic.NewUDPTracker(*localConn, statistic.DefaultManager, ct.Metadata(), ct.Rule())

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer func() {
			wg.Done()
			once.Do(cleanConn)
		}()
		err := copyUdpPacket(lc, ct.Conn())
		if err != nil {
			conn.PrintPacketError(err, fmt.Sprintf(log.FormatLog(log.UdpPrefix, "fail to output ,err: %v, remote address: %v"), err, ct.Metadata().RemoteAddress()))
		}
	}()
	go func() {
		defer func() {
			wg.Done()
			once.Do(cleanConn)
		}()
		err := copyUdpPacket(ct.Conn(), lc)
		if err != nil {
			conn.PrintPacketError(err, fmt.Sprintf(log.FormatLog(log.UdpPrefix, "fail to input ,err: %v, remote address: %v"), err, ct.Metadata().RemoteAddress()))
		}
	}()

	wg.Wait()
}

func handleDnsConn(ct conn.UdpConnContext) {
	var err error

	remoteConn := conn.CopyableReaderWriterConn{PacketConn: ct.Conn()}

	defer func() {
		ct.Wg().Done()
		err := closeConn(remoteConn)
		if err != nil {
			log.Debugln(log.FormatLog(log.DnsPrefix, "fail to close remote conn,err: %v"), err)
		}
	}()

	data := pool.NewBytes(pool.BufSize)
	defer pool.FreeBytes(data)
	_, addr, err := remoteConn.ReadFrom(data)
	if err != nil {
		log.Debugln(log.FormatLog(log.DnsPrefix, "fail to read dns message: err: %v"), err)
		return
	}
	dnsMessage := new(D.Msg)
	err = dnsMessage.Unpack(data)
	if err != nil {
		log.Warnln(log.FormatLog(log.DnsPrefix, "fail to unpack dns message: err: %v"), err)
		return
	}
	res, err := dns.Handle(dnsMessage, ct.Metadata())
	if err != nil {
		log.Warnln(log.FormatLog(log.DnsPrefix, "fail to handle dns message: err: %v"), err)
		return
	}
	resData, err := res.Pack()
	if err != nil {
		log.Warnln(log.FormatLog(log.DnsPrefix, "fail to pack dns message: err: %v"), err)
		return
	}
	_, err = remoteConn.WriteTo(resData, addr)
	if err != nil {
		log.Debugln(log.FormatLog(log.DnsPrefix, "fail to write back dns message: err: %v"), err)
	}
}

// processUDP starts a loop to handle udp packet
func processUDP() {
	for c := range udpQueue {
		if conn.GetIsDNSConn(c.Metadata()) {
			go handleDnsConn(c)
		} else {
			go handleUdpConn(c)
		}

	}
}
