package dns

import (
	"github.com/igoogolx/itun2socks/internal/constants"
	"go.uber.org/atomic"
)

type StatisticItem struct {
	Success *atomic.Int32 `json:"success"`
	Fail    *atomic.Int32 `json:"fail"`
}

type Statistic struct {
	Proxy  StatisticItem `json:"proxy"`
	Direct StatisticItem `json:"direct"`
}

var static = Statistic{
	Proxy: StatisticItem{
		Success: atomic.NewInt32(0),
		Fail:    atomic.NewInt32(0),
	},
	Direct: StatisticItem{
		Success: atomic.NewInt32(0),
		Fail:    atomic.NewInt32(0),
	},
}

func countSuccessQuery(policy constants.Policy) {
	if policy == constants.PolicyDirect {
		static.Direct.Success.Inc()
	} else if policy == constants.PolicyProxy {
		static.Proxy.Success.Inc()
	}
}

func countFailQuery(policy constants.Policy) {
	if policy == constants.PolicyDirect {
		static.Direct.Fail.Inc()
	} else if policy == constants.PolicyProxy {
		static.Proxy.Fail.Inc()
	}
}

func GetStatistic() Statistic {
	return static
}
