package tunnel

import (
	"fmt"
	"github.com/igoogolx/itun2socks/internal/conn"
	"github.com/igoogolx/itun2socks/internal/constants"
	"github.com/igoogolx/itun2socks/internal/dns"
	"github.com/igoogolx/itun2socks/internal/tunnel/statistic"
	"github.com/igoogolx/itun2socks/pkg/log"
	"github.com/igoogolx/itun2socks/pkg/network_iface"
	"github.com/igoogolx/itun2socks/pkg/pool"
	E "github.com/sagernet/sing/common/exceptions"
	"sync"
	"time"
)

func ShouldIgnorePacketError(err error) bool {
	// ignore simple error
	if E.IsTimeout(err) || E.IsClosed(err) || E.IsCanceled(err) {
		return true
	}
	return false
}

var (
	udpQueue = make(chan conn.UdpConnContext, 1024)
)

func UdpQueue() chan conn.UdpConnContext {
	return udpQueue
}

func copyUdpPacket(lc conn.UdpConn, rc conn.UdpConn) error {
	receivedBuf := pool.NewBytes(pool.BufSize)
	defer pool.FreeBytes(receivedBuf)
	for {

		err := rc.SetReadDeadline(time.Now().Add(5 * time.Second))
		if err != nil {
			return fmt.Errorf("fail to set udp conn read deadline: %v", err)
		}

		n, addr, err := rc.ReadFrom(receivedBuf)
		if ShouldIgnorePacketError(err) {
			log.Debugln("ignore packet read from err: %v", err)
			return nil
		}
		if err != nil {
			return fmt.Errorf("fail to read udp from rc:%v", err)
		}

		err = lc.SetWriteDeadline(time.Now().Add(5 * time.Second))
		if err != nil {
			return fmt.Errorf("fail to set udp conn write deadline: %v", err)
		}
		_, err = lc.WriteTo(receivedBuf[:n], addr)
		if ShouldIgnorePacketError(err) {
			return nil
		}
		if err != nil {
			return fmt.Errorf("fail to write udp to lc:%v", err)
		}
	}

}

func handleUdpConn(ct conn.UdpConnContext) {
	log.Debugln(log.FormatLog(log.UdpPrefix, "handle udp conn, remote address: %v"), ct.Metadata().RemoteAddress())
	defer func() {
		err := closeConn(ct.Conn())
		ct.Wg().Done()
		if err != nil {
			log.Warnln(log.FormatLog(log.UdpPrefix, "fail to close remote conn,err: %v"), err)
		}
		log.Debugln(log.FormatLog(log.UdpPrefix, "close remote conn: %v"), ct.Metadata().RemoteAddress())
	}()
	var lc conn.UdpConn
	var err error

	//only tun proxy
	if ct.Metadata().DstPort.String() == constants.DnsPort {
		err = dns.HandleDnsConn(ct.Conn())
		if err != nil {
			log.Warnln(log.FormatLog(log.UdpPrefix, "fail to handle dns conn, err: %v, remote address: %v"), err, ct.Metadata().RemoteAddress())
		}
		return
	} else {
		localConn, err := conn.NewUdpConn(ct.Ctx(), ct.Metadata(), ct.Rule(), network_iface.GetDefaultInterfaceName())
		if err != nil {
			log.Warnln(log.FormatLog(log.UdpPrefix, "fail to get udp conn, err: %v, remote address: %v"), err, ct.Metadata().RemoteAddress())
			return
		}
		lc = statistic.NewUDPTracker(localConn, statistic.DefaultManager, ct.Metadata(), ct.Rule())
	}

	defer func() {
		err = closeConn(lc)
		if err != nil {
			log.Warnln(log.FormatLog(log.UdpPrefix, "fail to close local conn,err: %v"), err)
		} else {
			log.Infoln(log.FormatLog(log.UdpPrefix, "close local conn, remote address: %v"), ct.Metadata().RemoteAddress())
		}
	}()

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		err := copyUdpPacket(lc, ct.Conn())
		if err != nil {
			log.Warnln(log.FormatLog(log.UdpPrefix, "fail to handle output ,err: %v, remote address: %v"), err, ct.Metadata().RemoteAddress())
		}
	}()
	go func() {
		defer wg.Done()
		err := copyUdpPacket(ct.Conn(), lc)
		if err != nil {
			log.Warnln(log.FormatLog(log.UdpPrefix, "fail to handle input ,err: %v, remote address: %v"), err, ct.Metadata().RemoteAddress())
		}
	}()

	wg.Wait()
}

// processUDP starts a loop to handle udp packet
func processUDP() {
	for c := range udpQueue {
		go handleUdpConn(c)
	}
}
