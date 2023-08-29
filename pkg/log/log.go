package log

import (
	"github.com/igoogolx/itun2socks/internal/constants"
	log "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

type output struct {
	logger *lumberjack.Logger
}

func (o *output) Write(p []byte) (int, error) {
	_, err := o.logger.Write(p)
	if err != nil {
		return 0, err
	}
	return os.Stdout.Write(p)
}

func init() {
	log.SetOutput(&output{
		&lumberjack.Logger{
			Filename:   constants.LogFile,
			MaxSize:    5, // megabytes
			MaxBackups: 1,
			MaxAge:     5,    //days
			Compress:   true, // disabled by default
		},
	})
}
