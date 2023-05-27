package realdb

import "time"

type Config struct {
	DSN         string        `toml:"dsn,omitempty"`
	MaxLifetime time.Duration `toml:"max_lifetime,omitempty"`
}

func DefaultConfig() *Config {
	return &Config{
		MaxLifetime: 24 * time.Hour,
		DSN:         "root:@tcp(localhost:4000)/",
	}
}

func (c *Config) SetDSN(dsn string) {
	c.DSN = dsn
}
