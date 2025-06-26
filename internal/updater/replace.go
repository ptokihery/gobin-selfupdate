package updater

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func ReplaceBinary(newBinary string) error {
	current, err := os.Executable()
	if err != nil {
		return fmt.Errorf("unable to get executable path: %w", err)
	}

	switch runtime.GOOS {
	case "windows":
		bat := fmt.Sprintf(`@echo off
timeout /t 1 > nul
move /Y "%s" "%s"
start "" "%s"
`, newBinary, current, current)
		batFile := filepath.Join(os.TempDir(), "update.bat")
		if err := os.WriteFile(batFile, []byte(bat), 0644); err != nil {
			return err
		}
		if err := exec.Command("cmd", "/C", batFile).Start(); err != nil {
			return err
		}

	default:
		sh := fmt.Sprintf(`#!/bin/sh
sleep 1
mv "%s" "%s"
chmod +x "%s"
"%s" &
`, newBinary, current, current, current)
		shFile := filepath.Join(os.TempDir(), "update.sh")
		if err := os.WriteFile(shFile, []byte(sh), 0755); err != nil {
			return err
		}
		if err := exec.Command("sh", shFile).Start(); err != nil {
			return err
		}
	}

	os.Exit(0)
	return nil
}
