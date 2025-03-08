//go:build windows

package service

import (
	"github.com/igoogolx/itun2socks/pkg/log"
	"github.com/kardianos/service"
)

type program struct{}

func (p *program) Start(s service.Service) error {
	// Start should not block. Do the actual work async.
	go p.run()
	return nil
}
func (p *program) run() {
	// Do work here
}
func (p *program) Stop(s service.Service) error {
	// Stop should not block. Return with a few seconds.
	return nil
}

var svcConfig = &service.Config{
	Name:        "LuxCoreGuardService",
	DisplayName: "Lux Core Service",
	Description: "Run the lux core as admin on default",
}
var s service.Service
var isInService = !service.Interactive()

func init() {
	if !isInService {
		return
	}
	prg := &program{}
	var err error
	s, err = service.New(prg, svcConfig)
	if err != nil {
		log.Errorln("fail to create service: %v", err)
	}
}

func Run() error {
	if !isInService {
		return nil
	}
	return s.Run()
}

func Install() error {
	if !isInService {
		return nil
	}
	return s.Install()
}

func Uninstall() error {
	if !isInService {
		return nil
	}
	return s.Uninstall()
}

func Restart() error {
	if !isInService {
		return nil
	}
	return s.Restart()
}

func Interactive() bool {
	return false
}
