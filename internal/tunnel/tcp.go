package tunnel

import (
	"github.com/igoogolx/itun2socks/internal/conn"
	"github.com/igoogolx/itun2socks/internal/tunnel/statistic"
	"github.com/igoogolx/itun2socks/pkg/log"
	"github.com/igoogolx/itun2socks/pkg/network_iface"
	"github.com/igoogolx/itun2socks/pkg/pool"
	"io"
	"net"
	"sync"
	"time"
)

var (
	tcpQueue   = make(chan conn.TcpConnContext, 1024)
	tcpTimeout = 5 * time.Minute
)

func TcpQueue() chan conn.TcpConnContext {
	return tcpQueue
}

func handleTCPConn(ct conn.TcpConnContext) {
	remoteConn, err := conn.NewTcpConn(ct.Ctx(), ct.Metadata(), ct.Rule(), network_iface.GetDefaultInterfaceName())
	defer func() {
		ct.Wg().Done()
		if err := closeConn(ct.Conn()); err != nil {
			log.Debugln(log.FormatLog(log.TcpPrefix, "fail to close local tcp conn,err: %v"), err)
		}
		if err := closeConn(remoteConn); err != nil {
			log.Debugln(log.FormatLog(log.TcpPrefix, "fail to close remote tcp conn, err: %v"), err)
		}
	}()
	if err != nil {
		log.Warnln(log.FormatLog(log.TcpPrefix, "fail to get tcp conn, err: %v, rule: %v, remote address: %v"), err, ct.Rule().GetPolicy(), ct.Metadata().RemoteAddress())
		return
	}
	remoteConn = statistic.NewTCPTracker(remoteConn, statistic.DefaultManager, ct.Metadata(), ct.Rule())

	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		if err := copyPacket(ct.Conn(), remoteConn); err != nil {

			log.Warnln(log.FormatLog(log.TcpPrefix, "fail to input: %v, remote address: %v"), err, ct.Metadata().RemoteAddress())
		}
	}()
	go func() {
		defer wg.Done()
		if err := copyPacket(remoteConn, ct.Conn()); err != nil {
			log.Warnln(log.FormatLog(log.TcpPrefix, "fail to output: %v, remote address: %v"), err, ct.Metadata().RemoteAddress())
		}
	}()
	wg.Wait()
}

func copyPacket(lc net.Conn, rc net.Conn) error {
	buf := pool.NewBytes(pool.BufSize)
	defer pool.FreeBytes(buf)
	err := lc.SetDeadline(time.Now().Add(tcpTimeout))
	if err != nil {
		return err
	}
	err = rc.SetDeadline(time.Now().Add(tcpTimeout))
	if err != nil {
		return err
	}
	_, err = io.CopyBuffer(lc, rc, buf)
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
