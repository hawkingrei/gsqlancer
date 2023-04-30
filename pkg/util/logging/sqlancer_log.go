package logging

type StdLogger struct {
	cfg *LogConfig
}

func NewStdLogger(cfg *LogConfig) *StdLogger {
	return &StdLogger{
		cfg,
	}
}
