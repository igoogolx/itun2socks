package statistic

import (
	"github.com/igoogolx/itun2socks/internal/constants"
	"github.com/igoogolx/itun2socks/pkg/log"
	"sync"
	"time"

	"go.uber.org/atomic"
)

var DefaultManager *Manager

func init() {
	DefaultManager = &Manager{
		proxy: &Statistic{
			uploadTemp:    atomic.NewInt64(0),
			downloadTemp:  atomic.NewInt64(0),
			uploadBlip:    atomic.NewInt64(0),
			downloadBlip:  atomic.NewInt64(0),
			uploadTotal:   atomic.NewInt64(0),
			downloadTotal: atomic.NewInt64(0),
		},
		direct: &Statistic{
			uploadTemp:    atomic.NewInt64(0),
			downloadTemp:  atomic.NewInt64(0),
			uploadBlip:    atomic.NewInt64(0),
			downloadBlip:  atomic.NewInt64(0),
			uploadTotal:   atomic.NewInt64(0),
			downloadTotal: atomic.NewInt64(0),
		},
	}

	go DefaultManager.handle(constants.DistributionProxy)
	go DefaultManager.handle(constants.DistributionBypass)
}

type Statistic struct {
	uploadTemp    *atomic.Int64
	downloadTemp  *atomic.Int64
	uploadBlip    *atomic.Int64
	downloadBlip  *atomic.Int64
	uploadTotal   *atomic.Int64
	downloadTotal *atomic.Int64
}

type Manager struct {
	connections sync.Map
	proxy       *Statistic
	direct      *Statistic
}

func (m *Manager) Join(c tracker) {
	m.connections.Store(c.ID(), c)
}

func (m *Manager) Leave(c tracker) {
	m.connections.Delete(c.ID())
}

func (m *Manager) getStatistic(rule constants.RuleType) *Statistic {
	if rule == constants.DistributionProxy {
		return m.proxy
	}
	if rule == constants.DistributionBypass {
		return m.direct
	}
	return nil
}

func (m *Manager) PushUploaded(size int64, rule constants.RuleType) {
	s := m.getStatistic(rule)
	if s != nil {
		s.uploadTemp.Add(size)
		s.uploadTotal.Add(size)
	}
}

func (m *Manager) PushDownloaded(size int64, rule constants.RuleType) {
	s := m.getStatistic(rule)
	if s != nil {
		s.downloadTemp.Add(size)
		s.downloadTotal.Add(size)
	}
}

func (m *Manager) Now(rule constants.RuleType) (up int64, down int64) {
	s := m.getStatistic(rule)
	if s != nil {
		return s.uploadBlip.Load(), s.downloadBlip.Load()
	}
	return 0, 0
}

func (m *Manager) Connections() []tracker {
	var connections []tracker
	connections = []tracker{}
	m.connections.Range(func(key, value interface{}) bool {
		connections = append(connections, value.(tracker))
		return true
	})

	return connections
}

func (m *Manager) CloseAllConnections() {
	connections := m.Connections()
	for _, c := range connections {
		err := c.Close()
		if err != nil {
			log.Debugln(log.FormatLog(log.ExecutorPrefix, "fail to close connection, err: %v"), err)
		}
	}
}

func (m *Manager) GetTotal() *Total {

	return &Total{
		SnapshotStatistic{
			m.proxy.downloadTotal.Load(),
			m.proxy.uploadTotal.Load(),
		},
		SnapshotStatistic{
			m.direct.downloadTotal.Load(),
			m.direct.uploadTotal.Load(),
		},
	}
}

func (m *Manager) ResetStatistic(rule constants.RuleType) {
	s := m.getStatistic(rule)
	if s != nil {
		s.uploadTemp.Store(0)
		s.uploadBlip.Store(0)
		s.uploadTotal.Store(0)
		s.downloadTemp.Store(0)
		s.downloadBlip.Store(0)
		s.downloadTotal.Store(0)
	}
}

func (m *Manager) handle(rule constants.RuleType) {
	s := m.getStatistic(rule)
	if s == nil {
		return
	}
	ticker := time.NewTicker(time.Second)
	for range ticker.C {
		s.uploadBlip.Store(s.uploadTemp.Load())
		s.uploadTemp.Store(0)
		s.downloadBlip.Store(s.downloadTemp.Load())
		s.downloadTemp.Store(0)
	}
}

type SnapshotStatistic struct {
	DownloadTotal int64 `json:"download"`
	UploadTotal   int64 `json:"upload"`
}

type Total struct {
	Proxy  SnapshotStatistic `json:"proxy"`
	Direct SnapshotStatistic `json:"direct"`
}
