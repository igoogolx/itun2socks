package api

import (
	"fmt"
	"github.com/igoogolx/itun2socks/api/routes"
	"github.com/igoogolx/itun2socks/internal/constants"
	"github.com/igoogolx/itun2socks/pkg/log"
	"net"
	"os"
	"strconv"
)

func getFreePort() (int, error) {
	var a *net.TCPAddr
	var err error
	if a, err = net.ResolveTCPAddr("tcp", "localhost:0"); err == nil {
		var l *net.TCPListener
		if l, err = net.ListenTCP("tcp", a); err == nil {
			defer func(l *net.TCPListener) {
				err := l.Close()
				if err != nil {
					log.Errorln("fail to close listener:%v", err)
				}
			}(l)
			return l.Addr().(*net.TCPAddr).Port, nil
		}
	}
	return 0, fmt.Errorf("not found: %v", err)
}

func Start(port int, secret string) {
	go func() {
		var availablePort = port
		if availablePort == 0 {
			var err error
			availablePort, err = getFreePort()
			if err != nil {
				log.Errorln("fail to get free port:%v", err)
				os.Exit(1)
			}
		}

		hubAddress := fmt.Sprintf("http://localhost:%v", availablePort)
		if len(secret) != 0 {
			hubAddress = fmt.Sprintf("%v?token=%v", hubAddress, secret)
		}
		constants.SetHubAddress(hubAddress)
		log.Infoln("%s", log.FormatLog(log.ExecutorPrefix, hubAddress))
		err := routes.Start("localhost:"+strconv.Itoa(availablePort), secret)
		if err != nil {
			log.Fatalln(log.FormatLog(log.ExecutorPrefix, "fail to start hub: %v"), err)
		}
	}()
}
