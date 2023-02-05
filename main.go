package main

import (
	"flag"
	"fmt"
	"github.com/Dreamacro/clash/log"
	_ "github.com/igoogolx/itun2socks/components/debug"
	"github.com/igoogolx/itun2socks/components/is-elevated"
	_ "github.com/igoogolx/itun2socks/components/log"
	configurationTypes "github.com/igoogolx/itun2socks/configuration/configuration-types"
	"github.com/igoogolx/itun2socks/constants"
	"github.com/igoogolx/itun2socks/hub"
	"github.com/igoogolx/itun2socks/manager"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	version := flag.Bool("version", false, "Print version")
	port := flag.Int("port", constants.DefaultHubPort, "Running port, default:9000")
	config := flag.String("config", constants.DbFileName, "Config file path, default: config.json")
	checkElevated := flag.Bool("check_elevated", true, "Check whether it's run as the admin, default: true")
	flag.Parse()
	if *version {
		fmt.Println(constants.Version)
		os.Exit(0)
	}
	if *checkElevated {
		if !is_elevated.Get() {
			log.Fatalln("Please run as administrator or root")
			return
		}
	}
	configurationTypes.ConfigFilePath.Store(*config)
	hub.Start(*port)
	defer func() {
		if p := recover(); p != nil {
			log.Errorln("internal error: %v", p)
		}
		err := manager.Close()
		if err != nil {
			log.Errorln("fail to close executor:%v", err)
		}
	}()
	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGHUP)
	<-osSignals
}
