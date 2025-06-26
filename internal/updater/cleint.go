package updater

import "context"

type Client interface {
	FetchManifest(ctx context.Context) (*Manifest, error)
	DownloadBinary(ctx context.Context, keyOrURL, dest string) error
}
