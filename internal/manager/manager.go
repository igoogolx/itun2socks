package manager

import (
	"errors"
	"github.com/igoogolx/itun2socks/internal/executor"
	"github.com/igoogolx/itun2socks/pkg/log"
	"sync"
)

var (
	client executor.Client
	mux    sync.Mutex
)

func Start() error {
	mux.Lock()
	defer mux.Unlock()
	var err error
	defer func() {
		if err != nil {
			client = nil
		}
	}()
	if GetIsStarted() {
		return errors.New("the client has started")
	}
	client, err = executor.New()
	if err != nil {
		return err
	}
	err = client.Start()
	if err != nil {
		log.Errorln(log.FormatLog(log.ExecutorPrefix, "fail to start the client: %v"), err)
		err = client.Close()
		if err != nil {
			log.Errorln(log.FormatLog(log.ExecutorPrefix, "fail to close the client: %v"), err)
		}
		return err
	}
	log.Infoln(log.FormatLog(log.ExecutorPrefix, "Started the client successfully"))
	return nil
}

func Close() error {
	mux.Lock()
	defer mux.Unlock()
	if client != nil {
		err := client.Close()
		client = nil
		if err != nil {
			return err
		}
	}
	log.Infoln(log.FormatLog(log.ExecutorPrefix, "Stopped the client successfully"))
	return nil
}

func GetIsStarted() bool {
	return client != nil
}

func RuntimeDetail(hubAddress string) (interface{}, error) {
	return client.RuntimeDetail(hubAddress)
}
