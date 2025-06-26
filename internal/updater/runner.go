package updater

import (
	"context"
	"fmt"

	"golang.org/x/mod/semver"
)

type Runner struct {
	Updater *Updater
}

func (r *Runner) Run(ctx context.Context, currentVersion string) error {
	manifest, err := r.Updater.CheckForUpdate(ctx)
	if err != nil {
		return err
	}

	cv := currentVersion
	mv := manifest.Version

	if !semver.IsValid(cv) && semver.IsValid("v"+cv) {
		cv = "v" + cv
	}
	if !semver.IsValid(mv) && semver.IsValid("v"+mv) {
		mv = "v" + mv
	}

	if !semver.IsValid(cv) || !semver.IsValid(mv) {
		return fmt.Errorf("invalid semver version: current=%q, manifest=%q", currentVersion, manifest.Version)
	}

	if semver.Compare(mv, cv) <= 0 {
		return nil // no update needed
	}

	return r.Updater.DownloadAndReplace(ctx, manifest, "/tmp/myapp_new")
}
