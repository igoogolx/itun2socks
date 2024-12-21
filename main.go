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
	"syscall"
)

var (
	homeDir string
	version bool
	port    int
	secret  string
)

var packageName = "itun2socks"

//go:generate go run scripts/generate.go
func main() {

	userConfigDir, _ := os.UserConfigDir()
	defaultHomeDir := filepath.Join(userConfigDir, packageName)

	flag.IntVar(&port, "port", constants.DefaultHubPort, "set running port")
	flag.StringVar(&homeDir, "home_dir", defaultHomeDir, "set configuration directory")
	flag.BoolVar(&version, "version", false, "print current version of itun2socks")
	flag.StringVar(&secret, "secret", "", "set secret")
	flag.Parse()

	if version {
		fmt.Printf("version: %v, build on: %v", constants.Version, constants.BuildTime)
		os.Exit(0)
	}

	if !filepath.IsAbs(homeDir) {
		currentDir, _ := os.Getwd()
		homeDir = filepath.Join(currentDir, homeDir)
	}
	_ = os.MkdirAll(homeDir, os.ModePerm)
	constants.Path.SetHomeDir(homeDir)
	log.InitLog()
	log.Infoln(log.FormatLog(log.InitPrefix, "using config: %v"), constants.Path.ConfigFilePath())
	configuration.SetConfigFilePath(constants.Path.ConfigFilePath())
	api.Start(port, secret)
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
