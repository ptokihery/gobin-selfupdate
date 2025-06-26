package config

import (
	internal "github.com/ptokihery/gobin-selfupdate/internal/config"
)

type Config = internal.Config

func Load(data []byte, key string) (*Config, error) {
	var cfg Config
	if err := internal.LoadEncryptedConfig(data, key, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}