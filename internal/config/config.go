package config

import (
	"encoding/hex"
	"encoding/json"
	"errors"
)

type Config struct {
	AWSRegion   string `json:"aws_region"`
	AccessKey   string `json:"access_key"`
	SecretKey   string `json:"secret_key"`
	Bucket      string `json:"bucket"`
	ManifestKey string `json:"manifest_key"`
	ManifestURL string `json:"manifest_url"`
	UseS3       bool   `json:"use_s3"`
	Interval       int   `json:"update_interval_minutes"`
}


func LoadEncryptedConfig(data []byte, keyHex string, cfg interface{}) error {
	key, err := hex.DecodeString(keyHex)
	if err != nil {
		return errors.New("invalid hex key")
	}
	if len(key) != 32 {
		return errors.New("key must be 32 bytes")
	}

	plaintext, err := decryptAES(data, key)
	if err != nil {
		return err
	}

	return json.Unmarshal(plaintext, cfg)
}
