package api

import (
	"fmt"
	"github.com/igoogolx/itun2socks/api/routes"
	"github.com/igoogolx/itun2socks/internal/constants"
	"github.com/igoogolx/itun2socks/pkg/log"
	"strconv"
)

func Start(port int, secret string) {
	go func() {
		hubAddress := fmt.Sprintf("http://localhost:%v", port)
		if len(secret) != 0 {
			hubAddress = fmt.Sprintf("%v?token=%v", hubAddress, secret)
		}
		constants.SetHubAddress(hubAddress)
		log.Infoln(log.FormatLog(log.ExecutorPrefix, hubAddress))
		err := routes.Start("localhost:"+strconv.Itoa(port), secret)
		if err != nil {
			log.Fatalln(log.FormatLog(log.ExecutorPrefix, "fail to start hub: %v"), err)
		}
	}()
}
