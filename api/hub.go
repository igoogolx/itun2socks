package api

import (
	"github.com/Dreamacro/clash/log"
	"github.com/igoogolx/itun2socks/api/routes"
	"strconv"
)

func Start(port int) {
	go func() {
		log.Infoln("hub is running on: http://localhost:%v", port)
		err := routes.Start("localhost:" + strconv.Itoa(port))
		if err != nil {
			log.Fatalln("fail to start hub: %v", err)
		}
	}()
}