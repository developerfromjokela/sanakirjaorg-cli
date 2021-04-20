package log

import (
	"github.com/withmandala/go-log"
	"os"
)

const (
	ERROR = 0
	WARN  = 1
	LOG   = 2
	DEBUG = 3
)

var LogLevel int

var logger = log.New(os.Stderr)

func Log(level int, message interface{}) {
	b := LogLevel == level || LogLevel > level
	logger.WithDebug()
	if b {
		switch level {
		case WARN:
			logger.Warn(message)
			break
		case ERROR:
			logger.Error(message)
			break
		case LOG:
			logger.Info(message)
			break
		case DEBUG:
			logger.Debug(message)
			break
		}
	}
}
