package config

import "github.com/hawkingrei/gsqlancer/pkg/connection"

type Config struct {
	EnablePartition bool              `json:"enable_partition,omitempty",default:"true"`
	db              connection.Config `json:"db"`
}

func DefaultConfig() *Config {
	return &Config{
		EnablePartition: true,
	}
}
