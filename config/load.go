package config

import (
	"embed"

	internal "github.com/ptokihery/gobin-selfupdate/internal/config"
)

type Config = internal.Config

func Load(fs embed.FS, file string, key []byte) (*Config, error) {
	var cfg Config
	if err := internal.LoadEncryptedConfig(fs, file, key, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
