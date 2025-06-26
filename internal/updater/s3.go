package updater

import (
	"context"
	"encoding/json"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Client struct {
	Client      *s3.Client
	Bucket      string
	ManifestKey string
}

func (c *S3Client) FetchManifest(ctx context.Context) (*Manifest, error) {
	resp, err := c.Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &c.Bucket,
		Key:    &c.ManifestKey,
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var m Manifest
	err = json.NewDecoder(resp.Body).Decode(&m)
	return &m, err
}

func (c *S3Client) DownloadBinary(ctx context.Context, key, dest string) error {
	resp, err := c.Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &c.Bucket,
		Key:    &key,
	})
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	f, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.ReadFrom(resp.Body)
	return err
}
