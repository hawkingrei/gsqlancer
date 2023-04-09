package config

type Config struct {
	EnablePartition bool
}

func DefaultConfig() *Config {
	return &Config{
		EnablePartition: true,
	}
}
