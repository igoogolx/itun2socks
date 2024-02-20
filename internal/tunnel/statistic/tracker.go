package statistic

import (
	"github.com/igoogolx/itun2socks/internal/cfg/distribution"
	"github.com/igoogolx/itun2socks/internal/constants"
	"net"
	"time"

	C "github.com/igoogolx/clash/constant"

	"github.com/gofrs/uuid/v5"
	"go.uber.org/atomic"
)

type tracker interface {
	ID() string
	Close() error
}

type trackerInfo struct {
	UUID          uuid.UUID        `json:"id"`
	Metadata      *C.Metadata      `json:"metadata"`
	UploadTotal   *atomic.Int64    `json:"upload"`
	DownloadTotal *atomic.Int64    `json:"download"`
	Start         int64            `json:"start"`
	Rule          constants.Policy `json:"rule"`
	Domain        string           `json:"domain"`
	Process       string           `json:"process"`
}

type TcpTracker struct {
	net.Conn `json:"-"`
	*trackerInfo
	manager *Manager
}

func (tt *TcpTracker) ID() string {
	return tt.UUID.String()
}

func (tt *TcpTracker) Read(b []byte) (int, error) {
	n, err := tt.Conn.Read(b)
	download := int64(n)
	tt.manager.PushDownloaded(download, tt.Rule)
	tt.DownloadTotal.Add(download)
	return n, err
}

func (tt *TcpTracker) Write(b []byte) (int, error) {
	n, err := tt.Conn.Write(b)
	upload := int64(n)
	tt.manager.PushUploaded(upload, tt.Rule)
	tt.UploadTotal.Add(upload)
	return n, err
}

func (tt *TcpTracker) Close() error {
	tt.manager.Leave(tt)
	return tt.Conn.Close()
}

func NewTCPTracker(conn net.Conn, manager *Manager, metadata *C.Metadata, rule constants.Policy) *TcpTracker {
	uid, _ := uuid.NewV4()

	t := &TcpTracker{
		Conn:    conn,
		manager: manager,
		trackerInfo: &trackerInfo{
			UUID:          uid,
			Start:         time.Now().UnixNano() / int64(time.Millisecond),
			Metadata:      metadata,
			Rule:          rule,
			UploadTotal:   atomic.NewInt64(0),
			DownloadTotal: atomic.NewInt64(0),
			Process:       "unknown",
		},
	}
	if cachedDomain, _, ok := distribution.GetCachedDnsItem(metadata.DstIP.String()); ok {
		t.trackerInfo.Domain = cachedDomain
	} else if metadata.Host != "" {
		t.trackerInfo.Domain = metadata.Host
	} else {
		t.trackerInfo.Domain = "unknown"
	}
	manager.Join(t)
	return t
}

type UdpTracker struct {
	net.PacketConn `json:"-"`
	*trackerInfo
	manager *Manager
}

func (ut *UdpTracker) ID() string {
	return ut.UUID.String()
}

func (ut *UdpTracker) ReadFrom(b []byte) (int, net.Addr, error) {
	n, addr, err := ut.PacketConn.ReadFrom(b)
	download := int64(n)
	ut.manager.PushDownloaded(download, ut.Rule)
	ut.DownloadTotal.Add(download)
	return n, addr, err
}

func (ut *UdpTracker) WriteTo(b []byte, addr net.Addr) (int, error) {
	n, err := ut.PacketConn.WriteTo(b, addr)
	upload := int64(n)
	ut.manager.PushUploaded(upload, ut.Rule)
	ut.UploadTotal.Add(upload)
	return n, err
}

func (ut *UdpTracker) Close() error {
	ut.manager.Leave(ut)
	return ut.PacketConn.Close()
}

func NewUDPTracker(conn net.PacketConn, manager *Manager, metadata *C.Metadata, rule constants.Policy) *UdpTracker {
	uid, _ := uuid.NewV4()

	ut := &UdpTracker{
		PacketConn: conn,
		manager:    manager,
		trackerInfo: &trackerInfo{
			UUID:          uid,
			Start:         time.Now().UnixNano() / int64(time.Millisecond),
			Metadata:      metadata,
			Rule:          rule,
			UploadTotal:   atomic.NewInt64(0),
			DownloadTotal: atomic.NewInt64(0),
			Process:       "unknown",
		},
	}

	if cachedDomain, _, ok := distribution.GetCachedDnsItem(metadata.DstIP.String()); ok {
		ut.trackerInfo.Domain = cachedDomain
	} else if metadata.Host != "" {
		ut.trackerInfo.Domain = metadata.Host
	} else {
		ut.trackerInfo.Domain = "unknown"
	}

	manager.Join(ut)
	return ut
}
