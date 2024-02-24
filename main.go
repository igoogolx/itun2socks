package main

import (
	"flag"
	"fmt"
	"github.com/igoogolx/itun2socks/api"
	"github.com/igoogolx/itun2socks/internal/configuration"
	"github.com/igoogolx/itun2socks/internal/constants"
	"github.com/igoogolx/itun2socks/internal/manager"
	"github.com/igoogolx/itun2socks/pkg/log"
	"os"
	"os/signal"
	"path/filepath"
	"runtime/debug"
	"syscall"
)

var (
	homeDir string
	version bool
	port    int
)

//go:generate go run scripts/generate.go
func main() {
	debug.SetMemoryLimit(32 * 1024 * 1024)
	flag.BoolVar(&version, "version", false, "Print version")
	flag.IntVar(&port, "port", constants.DefaultHubPort, "Running port, default: ${user_config_dir}/itun2socks")
	flag.StringVar(&homeDir, "home_dir", "", "Config dir, default: current dir")
	flag.Parse()

	if homeDir != "" {
		if !filepath.IsAbs(homeDir) {
			configDir, _ := os.UserConfigDir()
			homeDir = filepath.Join(configDir, "itun2socks", homeDir)
		}
	} else {
		configDir, _ := os.UserConfigDir()
		homeDir = filepath.Join(configDir, "itun2socks")
	}
	_ = os.MkdirAll(homeDir, os.ModePerm)
	constants.Path.SetHomeDir(homeDir)
	log.InitLog()
	log.Infoln(log.FormatLog(log.InitPrefix, "using config: %v"), constants.Path.ConfigFilePath())
	configuration.SetConfigFilePath(constants.Path.ConfigFilePath())

	if version {
		fmt.Printf("version: %v, build on: %v", constants.Version, constants.BuildTime)
		os.Exit(0)
	}
	api.Start(port)
	defer func() {
		if p := recover(); p != nil {
			log.Errorln(log.FormatLog(log.InitPrefix, "internal error: %v"), p)
		}
		err := manager.Close()
		if err != nil {
			log.Errorln(log.FormatLog(log.InitPrefix, "fail to close client:%v"), err)
		}
	}()
	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGHUP)
	<-osSignals
}
