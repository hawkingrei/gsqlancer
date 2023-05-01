package logging

type LogConfig struct {
	StatusLogPath string `toml:"status-log-path"`
	SQLLogPath    string `toml:"sql-log-path"`
}

func (l *LogConfig) FileLogConfig() *FileLogConfig {
	return &FileLogConfig{
		Filename: l.SQLLogPath,
	}
}

// FileLogConfig serializes file log related config in toml/json.
type FileLogConfig struct {
	// Log filename, leave empty to disable file log.
	Filename string `toml:"filename" json:"filename"`
	// Max size for a single file, in MB.
	MaxSize int `toml:"max-size" json:"max-size"`
	// Max log keep days, default is never deleting.
	MaxDays int `toml:"max-days" json:"max-days"`
	// Maximum number of old log files to retain.
	MaxBackups int `toml:"max-backups" json:"max-backups"`
}
