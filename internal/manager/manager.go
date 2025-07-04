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
	var err, startErr, closeErr error
	defer func() {
		if err != nil || startErr != nil || closeErr != nil {
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
	startErr = client.Start()
	if startErr != nil {
		log.Errorln(log.FormatLog(log.ExecutorPrefix, "fail to start the client: %v"), startErr)
		closeErr = client.Close()
		if closeErr != nil {
			log.Errorln(log.FormatLog(log.ExecutorPrefix, "fail to close the client: %v"), closeErr)
		}
		return startErr
	}
	log.Infoln("%s", log.FormatLog(log.ExecutorPrefix, "started the client successfully"))
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
	log.Infoln("%s", log.FormatLog(log.ExecutorPrefix, "stopped the client successfully"))
	return nil
}

func GetIsStarted() bool {
	return client != nil
}

func RuntimeDetail(hubAddress string) (interface{}, error) {
	return client.RuntimeDetail(hubAddress)
}
