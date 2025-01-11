package dns

import "go.uber.org/atomic"

var failCount = atomic.NewInt32(0)
var successCount = atomic.NewInt32(0)

func countSuccessQuery() {
	successCount.Inc()
}

func countFailQuery() {
	failCount.Inc()
}

func GetSuccessQueryCount() int32 {
	return successCount.Load()
}

func GetFailQueryCount() int32 {
	return failCount.Load()
}
