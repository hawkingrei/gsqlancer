package config

import (
	"time"

	"github.com/hawkingrei/gsqlancer/pkg/connection"
)

type Config struct {
	enablePartition bool              `toml:"enable_partition,omitempty"`
	concurrency     int32             `toml:"concurrency,omitempty"`
	maxTestTime     time.Duration     `toml:"max_test_time,omitempty"`
	db              connection.Config `toml:"db"`
}

func DefaultConfig() *Config {
	return &Config{
		enablePartition: true,
		concurrency:     8,
		maxTestTime:     6 * time.Hour,
		db:              *connection.DefaultConfig(),
	}
}

func (c *Config) Concurrency() int32 {
	return c.concurrency
}

func (c *Config) EnablePartition() bool {
	return c.enablePartition
}

func (c *Config) DBConfig() *connection.Config {
	return &c.db
}

func (c *Config) MaxTestTime() time.Duration {
	return c.maxTestTime
}
