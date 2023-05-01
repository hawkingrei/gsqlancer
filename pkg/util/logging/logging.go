package logging

import (
	"sync/atomic"

	"go.uber.org/zap"
)

var globalLogger atomic.Pointer[Logger]

type Logger interface {
	StatusLog() *zap.Logger
	SQLLog() *zap.Logger
	Stop()
}

var _ Logger = &StdLogger{}

func SetGlobalLogger(log *Logger) {
	globalLogger.Store(log)
}

func StatusLog() *zap.Logger {
	logger := globalLogger.Load()
	return (*logger).StatusLog()
}

func SQLLOG() *zap.Logger {
	logger := globalLogger.Load()
	return (*logger).SQLLog()
}
