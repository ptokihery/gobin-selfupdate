package updater

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
)

type Updater struct {
	Client Client
}

func (u *Updater) CheckForUpdate(ctx context.Context) (*Manifest, error) {
	return u.Client.FetchManifest(ctx)
}

func (u *Updater) DownloadAndReplace(ctx context.Context, manifest *Manifest, dest string) error {
	err := u.Client.DownloadBinary(ctx, manifest.ObjectKey, dest)
	if err != nil {
		return fmt.Errorf("download failed: %w", err)
	}

	file, err := os.Open(dest)
	if err != nil {
		return err
	}
	defer file.Close()

	h := sha256.New()
	if _, err := io.Copy(h, file); err != nil {
		return err
	}

	actualSum := fmt.Sprintf("%x", h.Sum(nil))
	if actualSum != manifest.ChecksumSHA {
		return fmt.Errorf("checksum mismatch: expected %s, got %s", manifest.ChecksumSHA, actualSum)
	}

	return ReplaceBinary(dest)
}
