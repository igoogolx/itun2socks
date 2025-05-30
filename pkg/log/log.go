package log

import (
	cLog "github.com/Dreamacro/clash/log"
	"github.com/igoogolx/itun2socks/internal/constants"
	"github.com/igoogolx/itun2socks/pkg/is_elevated"
	log "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"runtime"
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

func InitLog() {
	filePath := constants.Path.LogFilePath()
	if runtime.GOOS == "darwin" && is_elevated.Get() {
		err := os.Chmod(filePath, 0644)
		if err != nil {
			log.Warnln("fail to chang log file permissions:", err)
		}
	}
	log.SetFormatter(&log.JSONFormatter{})
	log.StandardLogger()
	log.SetOutput(&output{
		&lumberjack.Logger{
			Filename:   filePath,
			MaxSize:    1, // megabytes
			MaxBackups: 1,
		},
	})
}

type Prefix int64

const (
	InitPrefix Prefix = iota
	ConfigurationPrefix
	HubPrefix
	ExecutorPrefix
	DnsPrefix
	RulePrefix
	TcpPrefix
	TunPrefix
	UdpPrefix
	CachePrefix
)

func FormatLog(prefix Prefix, content string) string {
	switch prefix {
	case InitPrefix:
		return "[INIT], " + content
	case ConfigurationPrefix:
		return "[Configuration], " + content
	case HubPrefix:
		return "[HUB], " + content
	case ExecutorPrefix:
		return "[EXECUTOR], " + content
	case DnsPrefix:
		return "[DNS], " + content
	case RulePrefix:
		return "[RULE], " + content
	case TcpPrefix:
		return "[TCP], " + content
	case TunPrefix:
		return "[TUN], " + content
	case UdpPrefix:
		return "[UDP], " + content
	case CachePrefix:
		return "[CACHE], " + content
	}
	return "[UNKNOWN]"
}

func Infoln(format string, v ...any) {
	cLog.Infoln(format, v...)
}

func Warnln(format string, v ...any) {
	cLog.Warnln(format, v...)
}

func Errorln(format string, v ...any) {
	cLog.Errorln(format, v...)
}

func Debugln(format string, v ...any) {
	cLog.Debugln(format, v...)
}

func Fatalln(format string, v ...any) {
	cLog.Fatalln(format, v...)
}
