package config

import (
	"time"

	"github.com/hawkingrei/gsqlancer/pkg/connection/realdb"
	"github.com/hawkingrei/gsqlancer/pkg/util/logging"
)

type Config struct {
	log                   *logging.LogConfig `toml:"log"`
	db                    realdb.Config      `toml:"db"`
	maxTestTime           time.Duration      `toml:"max_test_time"`
	concurrency           int32              `toml:"concurrency"`
	enablePartition       bool               `toml:"enable_partition"`
	enableTiflashReplicas bool               `toml:"enable_tiflash_replicas"`
	selectDepth           int                `toml:"select_depth"`

	ReportPath           string `toml:"report_path"`
	EnablePQSApproach    bool   `toml:"enable_pqs_approach"`
	EnableNoRECApproach  bool   `toml:"enable_no_rec_approach"`
	EnableTLPApproach    bool   `toml:"enable_tlp_approach"`
	ViewCount            int    `toml:"view_count"`
	EnableLeftRightJoin  bool   `toml:"enable_left_right_join"`
	IsInUpdateDeleteStmt bool   `toml:"is_in_update_delete_stmt"`
	IsInExprIndex        bool   `toml:"is_in_expr_index"`
	Depth                int    `toml:"depth"`
	Hint                 bool   `toml:"hint"`
}

func DefaultConfig() *Config {
	return &Config{
		log: &logging.LogConfig{
			StatusLogPath: "./status.log",
			SQLLogPath:    "./sql.log",
		},
		enablePartition:   true,
		concurrency:       8,
		maxTestTime:       6 * time.Hour,
		db:                *realdb.DefaultConfig(),
		EnablePQSApproach: true,
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
	return c.log
}
