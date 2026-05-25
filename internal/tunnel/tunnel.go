package tunnel

import (
	"runtime"
)

func process() {
	numUDPWorkers := max(runtime.NumCPU(), 4)
	for i := 0; i < numUDPWorkers; i++ {
		go processUDP()
	}
	go processTCP()
}

func init() {
	go process()
}
