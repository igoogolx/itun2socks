package main

import (
	"flag"
	"fmt"
	"github.com/Dreamacro/clash/log"
	_ "github.com/igoogolx/itun2socks/components/debug"
	"github.com/igoogolx/itun2socks/components/is-elevated"
	_ "github.com/igoogolx/itun2socks/components/log"
	"github.com/igoogolx/itun2socks/configuration"
	"github.com/igoogolx/itun2socks/constants"
	"github.com/igoogolx/itun2socks/hub"
	"github.com/igoogolx/itun2socks/manager"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
)

func main() {
	homeDir := *flag.String("home-dir", "", "Config dir, default: current dir")
	version := flag.Bool("version", false, "Print version")
	port := flag.Int("port", constants.DefaultHubPort, "Running port, default:9000")
	configFile := *flag.String("configFile", "", "Config file path, default: configFile.json")
	checkElevated := flag.Bool("check_elevated", true, "Check whether it's run as the admin, default: true")
	flag.Parse()

	if homeDir != "" {
		if !filepath.IsAbs(homeDir) {
			currentDir, _ := os.Getwd()
			homeDir = filepath.Join(currentDir, homeDir)
		}
		constants.Path.SetHomeDir(homeDir)
	}

	if configFile != "" {
		if !filepath.IsAbs(configFile) {
			currentDir, _ := os.Getwd()
			configFile = filepath.Join(currentDir, configFile)
		}
		configuration.SetConfigFilePath(configFile)
	} else {
		configuration.SetConfigFilePath(constants.Path.ConfigFilePath())
	}

	if *version {
		fmt.Printf("version: %v, build on: %v", constants.Version, constants.BuildTime)
		os.Exit(0)
	}
	if *checkElevated {
		if !is_elevated.Get() {
			log.Fatalln("Please run as administrator or root")
			return
		}
	}
	hub.Start(*port, constants.Path.WebDir())
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
