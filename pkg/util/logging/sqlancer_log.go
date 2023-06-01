package logging

import (
	"os"

	"github.com/pingcap/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type StdLogger struct {
	cfg       *LogConfig
	statusLog *zap.Logger
	sqlLog    *zap.Logger
	bws       *zapcore.BufferedWriteSyncer
}

func NewStdLogger(cfg *LogConfig) (*StdLogger, error) {
	encoder, err := log.NewTextEncoder(&log.Config{})
	if err != nil {
		return nil, err
	}
	// info level enabler
	infoLevel := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level == zapcore.DebugLevel
	})

	// error and fatal level enabler
	errorFatalLevel := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level == zapcore.ErrorLevel || level == zapcore.FatalLevel
	})

	// write syncers
	stdoutSyncer := zapcore.Lock(os.Stdout)
	stderrSyncer := zapcore.Lock(os.Stderr)

	// tee core
	core := zapcore.NewTee(
		zapcore.NewCore(
			encoder,
			stdoutSyncer,
			infoLevel,
		),
		zapcore.NewCore(
			encoder,
			stderrSyncer,
			errorFatalLevel,
		),
	)
	statusLog := zap.New(core)
	lg, err := initFileLog(cfg.FileLogConfig())
	if err != nil {
		return nil, err
	}
	output := zapcore.AddSync(lg)
	bws := zapcore.BufferedWriteSyncer{
		WS: output,
	}
	sqlCore := zapcore.NewCore(encoder, &bws, zap.InfoLevel)
	sqlLog := zap.New(sqlCore)
	return &StdLogger{
		cfg,
		statusLog,
		sqlLog,
		&bws,
	}, nil
}

func (s *StdLogger) StatusLog() *zap.Logger {
	return s.statusLog
}

func (s *StdLogger) SQLLog() *zap.Logger {
	return s.sqlLog
}

func (s *StdLogger) Stop() {
	s.bws.Stop()
}
