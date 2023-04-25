package connection

import "time"

type Config struct {
	DSN         string
	MaxLifetime time.Duration // maximum amount of time a connection may be reused
}

func DefaultConfig() *Config {
	return &Config{
		MaxLifetime: 24 * time.Hour,
	}
}

func (c *Config) SetDSN(dsn string) {
	c.DSN = dsn
}
