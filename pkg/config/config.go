package config

import "github.com/hawkingrei/gsqlancer/pkg/connection"

type Config struct {
	EnablePartition bool              `toml:"enable_partition,omitempty",default:"true"`
	db              connection.Config `toml:"db"`
}

func DefaultConfig() *Config {
	return &Config{
		EnablePartition: true,
	}
}
