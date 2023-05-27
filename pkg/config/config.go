package config

import (
	"time"

	"github.com/hawkingrei/gsqlancer/pkg/connection/realdb"
	"github.com/hawkingrei/gsqlancer/pkg/util/logging"
)

type Config struct {
	log                   logging.LogConfig `toml:"log"`
	db                    realdb.Config     `toml:"db"`
	maxTestTime           time.Duration     `toml:"max_test_time,omitempty"`
	concurrency           int32             `toml:"concurrency,omitempty"`
	enablePartition       bool              `toml:"enable_partition,omitempty"`
	enableTiflashReplicas bool              `toml:"enable_tiflash_replicas,omitempty"`
	selectDepth           int               `toml:"select_depth,omitempty"`

	EnablePQSApproach    bool
	EnableNoRECApproach  bool
	EnableTLPApproach    bool
	ViewCount            int
	EnableLeftRightJoin  bool
	IsInUpdateDeleteStmt bool
	IsInExprIndex        bool
	Depth                int
	Hint                 bool
}

func DefaultConfig() *Config {
	return &Config{
		enablePartition: true,
		concurrency:     8,
		maxTestTime:     6 * time.Hour,
		db:              *realdb.DefaultConfig(),
	}
}

func (c *Config) Concurrency() int32 {
	return c.concurrency
}

func (c *Config) EnablePartition() bool {
	return c.enablePartition
}

func (c *Config) EnableTiflashReplicas() bool {
	return c.enableTiflashReplicas
}

func (c *Config) DBConfig() *realdb.Config {
	return &c.db
}

func (c *Config) MaxTestTime() time.Duration {
	return c.maxTestTime
}

func (c *Config) Log() *logging.LogConfig {
	return &c.log
}
