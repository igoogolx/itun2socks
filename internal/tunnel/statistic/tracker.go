package statistic

import (
	"github.com/Dreamacro/clash/component/process"
	"github.com/igoogolx/itun2socks/internal/cfg/distribution"
	"github.com/igoogolx/itun2socks/internal/constants"
	"net"
	"net/netip"
	"time"

	C "github.com/Dreamacro/clash/constant"

	"github.com/gofrs/uuid"
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
	Rule          constants.IpRule `json:"rule"`
	Domain        string           `json:"domain"`
	Process       string           `json:"process"`
}

type tcpTracker struct {
	net.Conn `json:"-"`
	*trackerInfo
	manager *Manager
}

func (tt *tcpTracker) ID() string {
	return tt.UUID.String()
}

func (tt *tcpTracker) Read(b []byte) (int, error) {
	n, err := tt.Conn.Read(b)
	download := int64(n)
	tt.manager.PushDownloaded(download, tt.Rule)
	tt.DownloadTotal.Add(download)
	return n, err
}

func (tt *tcpTracker) Write(b []byte) (int, error) {
	n, err := tt.Conn.Write(b)
	upload := int64(n)
	tt.manager.PushUploaded(upload, tt.Rule)
	tt.UploadTotal.Add(upload)
	return n, err
}

func (tt *tcpTracker) Close() error {
	tt.manager.Leave(tt)
	return tt.Conn.Close()
}

func NewTCPTracker(conn net.Conn, manager *Manager, metadata *C.Metadata, rule constants.IpRule) *tcpTracker {
	uid, _ := uuid.NewV4()

	t := &tcpTracker{
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
	go func() {
		srcIP, srcOk := netip.AddrFromSlice(metadata.SrcIP)
		destIp, destOk := netip.AddrFromSlice(metadata.DstIP)
		if srcOk && destOk {
			srcIP = srcIP.Unmap()
			destIp = destIp.Unmap()
			processName, err := process.FindProcessPath(metadata.NetWork.String(), netip.AddrPortFrom(srcIP, uint16(metadata.SrcPort)), netip.AddrPortFrom(destIp, uint16(metadata.DstPort)))
			if err == nil && len(processName) != 0 {
				t.Process = processName
			}
			manager.Join(t)
		}

	}()
	if cachedItem, ok := distribution.GetCachedDnsItem(metadata.DstIP.String()); ok {
		if ok {
			t.trackerInfo.Domain = cachedItem.Domain
		}
	} else {
		t.trackerInfo.Domain = "unknown"
	}
	return t
}

type udpTracker struct {
	net.PacketConn `json:"-"`
	*trackerInfo
	manager *Manager
}

func (ut *udpTracker) ID() string {
	return ut.UUID.String()
}

func (ut *udpTracker) ReadFrom(b []byte) (int, net.Addr, error) {
	n, addr, err := ut.PacketConn.ReadFrom(b)
	download := int64(n)
	ut.manager.PushDownloaded(download, ut.Rule)
	ut.DownloadTotal.Add(download)
	return n, addr, err
}

func (ut *udpTracker) WriteTo(b []byte, addr net.Addr) (int, error) {
	n, err := ut.PacketConn.WriteTo(b, addr)
	upload := int64(n)
	ut.manager.PushUploaded(upload, ut.Rule)
	ut.UploadTotal.Add(upload)
	return n, err
}

func (ut *udpTracker) Close() error {
	ut.manager.Leave(ut)
	return ut.PacketConn.Close()
}

func NewUDPTracker(conn net.PacketConn, manager *Manager, metadata *C.Metadata, rule constants.IpRule) *udpTracker {
	uid, _ := uuid.NewV4()

	ut := &udpTracker{
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

	if cacheItem, ok := distribution.GetCachedDnsItem(metadata.DstIP.String()); ok {
		ut.trackerInfo.Domain = cacheItem.Domain
	} else {
		ut.trackerInfo.Domain = "unknown"
	}

	go func() {
		srcIP, srcOk := netip.AddrFromSlice(metadata.SrcIP)
		destIp, destOk := netip.AddrFromSlice(metadata.DstIP)
		if srcOk && destOk {
			srcIP = srcIP.Unmap()
			destIp = destIp.Unmap()
			processName, err := process.FindProcessPath(metadata.NetWork.String(), netip.AddrPortFrom(srcIP, uint16(metadata.SrcPort)), netip.AddrPortFrom(destIp, uint16(metadata.DstPort)))
			if err == nil && len(processName) != 0 {
				ut.Process = processName
			}
			manager.Join(ut)
		}
	}()
	return ut
}
