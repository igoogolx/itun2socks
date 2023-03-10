package tunnel

import (
	"github.com/igoogolx/itun2socks/constants"
	"runtime"
)

type Matcher interface {
	GetRule(ip string) constants.IpRule
}

func process() {
	numUDPWorkers := 4
	if runtime.NumCPU() > numUDPWorkers {
		numUDPWorkers = runtime.NumCPU()
	}
	for i := 0; i < numUDPWorkers; i++ {
		go processUDP()
	}
	go processTCP()
}

func init() {
	go process()
}
