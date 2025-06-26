package updater

import (
	"context"
	"time"

	internalConfig "github.com/ptokihery/gobin-seltupdate/internal/config"
	internalSelfupdate "github.com/ptokihery/gobin-seltupdate/internal/selfupdate"
	internalUpdater "github.com/ptokihery/gobin-seltupdate/internal/updater"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Config = internalConfig.Config

type Updater struct {
	checker *internalSelfupdate.Checker
}

func StartUpdater(cfg *Config, currentVersion string, interval time.Duration) (*Updater, error) {
	var client internalUpdater.Client

	if cfg.UseS3 {
		awsCfg, err := config.LoadDefaultConfig(context.Background(),
			config.WithRegion(cfg.AWSRegion),
			config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(cfg.AccessKey, cfg.SecretKey, "")),
		)
		if err != nil {
			return nil, err
		}
		s3client := s3.NewFromConfig(awsCfg)
		client = &internalUpdater.S3Client{
			Client:      s3client,
			Bucket:      cfg.Bucket,
			ManifestKey: cfg.ManifestKey,
		}
	} else {
		client = &internalUpdater.HTTPClient{
			ManifestURL: cfg.ManifestURL,
		}
	}

	runner := &internalUpdater.Runner{
		Updater: &internalUpdater.Updater{Client: client},
	}

	checker := internalSelfupdate.NewChecker(runner, currentVersion, interval)
	checker.Start()

	return &Updater{checker: checker}, nil
}

func (u *Updater) Stop() {
	u.checker.Stop()
}
