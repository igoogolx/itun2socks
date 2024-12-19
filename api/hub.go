package api

import (
	"github.com/igoogolx/itun2socks/api/routes"
	"github.com/igoogolx/itun2socks/pkg/log"
	"strconv"
)

func Start(port int, secret string) {
	go func() {
		log.Infoln(log.FormatLog(log.ExecutorPrefix, "hub is running on: http://localhost:%v?token=%v"), port, secret)
		err := routes.Start("localhost:"+strconv.Itoa(port), secret)
		if err != nil {
			log.Fatalln(log.FormatLog(log.ExecutorPrefix, "fail to start hub: %v"), err)
		}
	}()
}
