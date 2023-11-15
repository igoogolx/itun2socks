package tunnel

import (
	"runtime"
)

func process() {
	numUDPWorkers := 4
	if runtime.NumCPU() > numUDPWorkers {
		numUDPWorkers = runtime.NumCPU()
	}
	for i := 0; i < numUDPWorkers; i++ {
		go processUDP()
		go processTCP()
	}
}

func init() {
	go process()
}
