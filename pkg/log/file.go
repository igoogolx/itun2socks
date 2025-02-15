package log

import (
	"github.com/igoogolx/itun2socks/internal/constants"
	"github.com/igoogolx/itun2socks/pkg/list"
	log "github.com/sirupsen/logrus"
	"os"
)

func ReadFile(maxLen int) ([]string, error) {
	var filePath = constants.Path.LogFilePath()
	logFileHandler, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	defer func(logFileHandler *os.File) {
		err := logFileHandler.Close()
		if err != nil {
			log.Warnln("close file error:", err)
		}
	}(logFileHandler)
	items, err := list.ParseFile(logFileHandler)
	if err != nil {
		return nil, err
	}
	if len(items) >= maxLen {
		return items[len(items)-maxLen:], nil
	}
	return items, nil
}
