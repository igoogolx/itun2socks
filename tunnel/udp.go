package tunnel

import (
	"fmt"
	"github.com/Dreamacro/clash/log"
	"github.com/igoogolx/itun2socks/components/network-iface"
	"github.com/igoogolx/itun2socks/components/pool"
	"github.com/igoogolx/itun2socks/conn"
	"github.com/igoogolx/itun2socks/constants"
	"github.com/igoogolx/itun2socks/dns"
	"github.com/igoogolx/itun2socks/tunnel/statistic"
	"io"
	"net"
	"sync"
	"time"
)

var (
	udpQueue = make(chan conn.UdpConnContext, 200)
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
			return fmt.Errorf("fail to set udp conn deadline: %v", err)
		}
		n, addr, err := rc.ReadFrom(receivedBuf)
		if ne, ok := err.(net.Error); ok && ne.Timeout() {
			log.Debugln("udp read io timeout")
			return nil /* ignore I/O timeout */
		} else if err == io.EOF {
			return nil
		} else if err != nil {
			return fmt.Errorf("fail to read udp data from local: %v", err)
		}
		_, err = lc.WriteTo(receivedBuf[:n], addr)
		if err != nil {
			return fmt.Errorf("fail to write udp to remote:%v", err)
		}
	}

}

func handleUdpConn(ct conn.UdpConnContext) {
	log.Debugln("handle udp conn, dst ip: %v, dst port: %v", ct.Metadata().DstIP.String(), ct.Metadata().DstPort)
	defer func() {
		ct.Wg().Done()
		err := closeConn(ct.Conn())
		if err != nil {
			log.Debugln("fail to close remote udp conn,err: %v", err)
		}
	}()
	var lc conn.UdpConn
	var err error
	if ct.Metadata().DstPort.String() == constants.DnsPort {
		lc = dns.NewConn()
	} else {
		localConn, err := conn.NewUdpConn(ct.Ctx(), ct.Metadata(), ct.Rule(), network_iface.GetDefaultInterfaceName())
		if err != nil {
			log.Warnln("fail to get udp conn, err: %v, target: %v", err, ct.Metadata().DstIP.String())
			return
		}
		lc = statistic.NewUDPTracker(localConn, statistic.DefaultManager, ct.Metadata(), ct.Rule())
	}

	defer func() {
		err = closeConn(lc)
		if err != nil {
			log.Warnln("fail to close remote local conn,err: %v", err)
		}
	}()

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		err := copyUdpPacket(lc, ct.Conn())
		if err != nil {
			log.Warnln("fail to handle udp output, err: %v", err)
		}
	}()
	go func() {
		defer wg.Done()
		err := copyUdpPacket(ct.Conn(), lc)
		if err != nil {
			log.Warnln("fail to handle udp input, err: %v", err)
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
