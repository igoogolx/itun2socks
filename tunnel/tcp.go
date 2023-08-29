package tunnel

import (
	"github.com/Dreamacro/clash/log"
	"github.com/igoogolx/itun2socks/conn"
	network_iface "github.com/igoogolx/itun2socks/pkg/network-iface"
	"github.com/igoogolx/itun2socks/pkg/pool"
	"github.com/igoogolx/itun2socks/tunnel/statistic"
	"io"
	"net"
	"sync"
)

var (
	tcpQueue = make(chan conn.TcpConnContext, 200)
)

func TcpQueue() chan conn.TcpConnContext {
	return tcpQueue
}

func handleTCPConn(ct conn.TcpConnContext) {
	remoteConn, err := conn.NewTcpConn(ct.Ctx(), ct.Metadata(), ct.Rule(), network_iface.GetDefaultInterfaceName())
	defer func() {
		ct.Wg().Done()
		if err := closeConn(ct.Conn()); err != nil {
			log.Debugln("failed to close local tcp conn,err: %v", err)
		}
		if err := closeConn(remoteConn); err != nil {
			log.Debugln("failed to close remote tcp conn,err: %v", err)
		}
	}()
	if err != nil {
		log.Warnln("failed to get tcp conn, err: %v, rule: %v, remote ip: %v", err, ct.Rule(), ct.Metadata().DstIP)
		return
	}
	remoteConn = statistic.NewTCPTracker(remoteConn, statistic.DefaultManager, ct.Metadata(), ct.Rule())

	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		if err := copyPacket(ct.Conn(), remoteConn); err != nil {
			log.Warnln("fail to input tcp: %v", err)
		}
	}()
	go func() {
		defer wg.Done()
		if err := copyPacket(remoteConn, ct.Conn()); err != nil {
			log.Warnln("fail to output tcp: %v", err)
		}
	}()
	wg.Wait()
}

func copyPacket(lc net.Conn, rc net.Conn) error {
	buf := pool.NewBytes(pool.BufSize)
	defer pool.FreeBytes(buf)
	_, err := io.CopyBuffer(lc, rc, buf)
	return err
}

func processTCP() {
	for c := range tcpQueue {
		go handleTCPConn(c)
	}
}

type CloseableConn interface {
	Close() error
}

func closeConn(conn CloseableConn) error {
	if conn != nil {
		return conn.Close()
	}
	return nil
}
