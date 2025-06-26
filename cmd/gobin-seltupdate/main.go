package main

import (
	"context"
	"embed"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	configs "github.com/ptokihery/gobin-selfupdate/internal/config"
	"github.com/ptokihery/gobin-selfupdate/internal/selfupdate"
	"github.com/ptokihery/gobin-selfupdate/internal/updater"
)

type Config struct {
	AWSRegion   string
	AccessKey   string
	SecretKey   string
	Bucket      string
	ManifestKey string
	ManifestURL string
	UseS3       bool
}

func createS3Client(cfg *configs.Config) *s3.Client {
	awsCfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(cfg.AWSRegion),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(cfg.AccessKey, cfg.SecretKey, "")),
	)
	if err != nil {
		log.Fatalf("Failed to create AWS client: %v", err)
	}
	return s3.NewFromConfig(awsCfg)
}

//go:embed config.json.enc
var encryptedFS embed.FS

var encryptionKey = []byte("1234567890abcdef1234567890abcdef")

func main() {
	// configs.EncryptFile("config.json", "config.json.enc", encryptionKey)
	var cfg configs.Config
	if err := configs.LoadEncryptedConfig(encryptedFS, "config.json.enc", encryptionKey, &cfg); err != nil {
		log.Fatalf("Failed to load encrypted config: %v", err)
	}

	var client updater.Client

	if cfg.UseS3 {
		s3client := createS3Client(&cfg)
		client = &updater.S3Client{
			Client:      s3client,
			Bucket:      cfg.Bucket,
			ManifestKey: cfg.ManifestKey,
		}
	} else {
		client = &updater.HTTPClient{
			ManifestURL: cfg.ManifestURL,
		}
	}

	myUpdater := &updater.Updater{Client: client}
	runner := &updater.Runner{Updater: myUpdater}

	const currentVersion = "1.0.0"

	checker := selfupdate.NewChecker(runner, currentVersion, 10*time.Minute)
	checker.Start()
	defer checker.Stop()

	log.Println("Application started. Press Ctrl+C to stop.")

	select {}
}
