package manager

import (
	"errors"
	executor2 "github.com/igoogolx/itun2socks/internal/executor"
	"github.com/igoogolx/itun2socks/pkg/log"
	"runtime/debug"
	"sync"
)

var (
	client *executor2.Client
	mux    sync.Mutex
)

func Start() error {
	mux.Lock()
	defer mux.Unlock()
	var err error
	if GetIsStarted() {
		return errors.New("the client has started")
	}
	client, err = executor2.New()
	if err != nil {
		return err
	}
	err = client.Start()
	if err != nil {
		log.Errorln(log.FormatLog(log.ExecutorPrefix, "fail to start the client: %v"), err)
		err := client.Close()
		if err != nil {
			log.Errorln(log.FormatLog(log.ExecutorPrefix, "fail to close the client: %v"), err)
		}
		client = nil
		return err
	}
	log.Infoln(log.FormatLog(log.ExecutorPrefix, "Started the client successfully"))
	return nil
}

func Close() error {
	mux.Lock()
	defer mux.Unlock()
	//Pay attention to this because it may lead to performance problem
	debug.FreeOSMemory()
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

func RuntimeDetail() (*executor2.Detail, error) {
	return client.RuntimeDetail()
}
