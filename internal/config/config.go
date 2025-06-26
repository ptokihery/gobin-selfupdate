package config

import (
	"embed"
	"encoding/json"
)




type Config struct {
	AWSRegion   string `json:"aws_region"`
	AccessKey   string `json:"access_key"`
	SecretKey   string `json:"secret_key"`
	Bucket      string `json:"bucket"`
	ManifestKey string `json:"manifest_key"`
	ManifestURL string `json:"manifest_url"`
	UseS3       bool   `json:"use_s3"`
	Interval       bool   `json:"update_interval_minutes"`
}

func LoadEncryptedConfig(fs embed.FS, fileName string, key []byte, out any) error {
	data, err := fs.ReadFile(fileName)
	if err != nil {
		return err
	}

	decrypted, err := decryptAES(data, key)
	if err != nil {
		return err
	}

	return json.Unmarshal(decrypted, out)
}
