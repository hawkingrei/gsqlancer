package logging

type LogConfig struct {
	StatusLogPath string `toml:"status-log-path"`
	SQLLogPath    string `toml:"sql-log-path"`
}
