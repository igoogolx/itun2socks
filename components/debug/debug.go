//go:build debug
// +build debug

package debug

import (
	"github.com/Dreamacro/clash/log"
	"net/http"
	_ "net/http/pprof"
)

func init() {
	log.SetLevel(log.DEBUG)
	go func() {
		err := http.ListenAndServe("localhost:6060", nil)
		if err != nil {
			log.Errorln("fail to listen on localhost:6060, err: %v", err)
		}
	}()

}
