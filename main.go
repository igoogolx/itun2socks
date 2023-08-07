package main

import (
	"flag"
	"fmt"
	"github.com/Dreamacro/clash/log"
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

var (
	homeDir string
	version bool
	port    int
)

func main() {
	flag.BoolVar(&version, "version", false, "Print version")
	flag.IntVar(&port, "port", constants.DefaultHubPort, "Running port, default:9000")
	flag.StringVar(&homeDir, "home_dir", "", "Config dir, default: current dir")
	flag.Parse()

	if homeDir != "" {
		if !filepath.IsAbs(homeDir) {
			configDir, _ := os.UserConfigDir()
			homeDir = filepath.Join(configDir, "itun2socks", homeDir)
		}
	} else {
		currentDir, _ := os.Getwd()
		homeDir = filepath.Join(currentDir)
	}
	_ = os.MkdirAll(homeDir, os.ModePerm)
	constants.Path.SetHomeDir(homeDir)

	fmt.Printf(constants.Path.ConfigFilePath())

	configuration.SetConfigFilePath(constants.Path.ConfigFilePath())

	if version {
		fmt.Printf("version: %v, build on: %v", constants.Version, constants.BuildTime)
		os.Exit(0)
	}
	hub.Start(port)
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
