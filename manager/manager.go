package manager

import (
	"errors"
	"github.com/Dreamacro/clash/log"
	"github.com/igoogolx/itun2socks/executor"
	"sync"
)

var (
	client *executor.Client
	mux    sync.Mutex
)

func Start() error {
	mux.Lock()
	defer mux.Unlock()
	var err error
	if GetIsStarted() {
		return errors.New("the client has started")
	}
	client, err = executor.New()
	if err != nil {
		return err
	}
	err = client.Start()
	if err != nil {
		log.Errorln("fail to start the client: %v", err)
		err := client.Close()
		if err != nil {
			log.Errorln("fail to close the client: %v, when there is error of starting", err)
		}
		client = nil
		return err
	}
	log.Infoln("Started the client")
	return nil
}

func Close() error {
	mux.Lock()
	defer mux.Unlock()
	if client != nil {
		err := client.Close()
		if err != nil {
			return err
		}
		client = nil
	}
	log.Infoln("Stopped the client")
	return nil
}

func GetIsStarted() bool {
	return client != nil
}
