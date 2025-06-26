package updater

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
)

type HTTPClient struct {
	ManifestURL string
}

func (c *HTTPClient) FetchManifest(ctx context.Context) (*Manifest, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.ManifestURL, nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var m Manifest
	err = json.NewDecoder(resp.Body).Decode(&m)
	return &m, err
}

func (c *HTTPClient) DownloadBinary(ctx context.Context, url, dest string) error {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	f, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, resp.Body)
	return err
}
