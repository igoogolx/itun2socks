package tunnel

import (
	"github.com/Dreamacro/clash/log"
	"github.com/igoogolx/itun2socks/components/pool"
	"github.com/igoogolx/itun2socks/conn"
	"github.com/igoogolx/itun2socks/global"
	statistic2 "github.com/igoogolx/itun2socks/tunnel/statistic"
	"io"
	"net"
	"sync"
)

var (
	tcpQueue = make(chan conn.TcpConnContext, 200)
)

func AddTcpConn(ct conn.TcpConnContext) {
	tcpQueue <- ct
}

func handleTCPConn(ct conn.TcpConnContext) {
	remoteConn, err := conn.NewTcpConn(ct.Ctx, ct.Metadata(), ct.Rule, global.GetDefaultInterfaceName())
	defer func() {
		ct.Wg.Done()
		if err := closeConn(ct.Conn()); err != nil {
			log.Debugln("failed to close local tcp conn,err: %v", err)
		}
		if err := closeConn(remoteConn); err != nil {
			log.Debugln("failed to close remote tcp conn,err: %v", err)
		}
	}()
	if err != nil {
		log.Warnln("failed to get tcp conn, err: %v", err)
		return
	}
	remoteConn = statistic2.NewTCPTracker(remoteConn, statistic2.DefaultManager, ct.Metadata(), ct.Rule)

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
