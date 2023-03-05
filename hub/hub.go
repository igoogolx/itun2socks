package hub

import (
	"github.com/Dreamacro/clash/log"
	"github.com/igoogolx/itun2socks/hub/routes"
	"strconv"
)

func Start(port int, webDir string) {
	go func() {
		log.Infoln("hub is running on: http://localhost:%v", port)
		err := routes.Start("localhost:"+strconv.Itoa(port), webDir)
		if err != nil {
			log.Fatalln("fail to start hub: %v", err)
		}
	}()
}
