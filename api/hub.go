package api

import (
	"github.com/igoogolx/itun2socks/api/routes"
	"github.com/igoogolx/itun2socks/pkg/log"
	"strconv"
)

func Start(port int) {
	go func() {
		log.Infoln(log.FormatLog(log.ExecutorPrefix, "hub is running on: http://localhost:%v"), port)
		err := routes.Start("localhost:" + strconv.Itoa(port))
		if err != nil {
			log.Fatalln(log.FormatLog(log.ExecutorPrefix, "fail to start hub: %v"), err)
		}
	}()
}
