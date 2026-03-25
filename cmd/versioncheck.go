package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// checkForUpdate prints a nudge if a newer version is available.
// Only checks once per day, stores last check time in ~/.dbx/.last_update_check.
func checkForUpdate() {
	if Version == "dev" {
		return
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return
	}

	checkFile := filepath.Join(home, ".dbx", ".last_update_check")

	// Skip if checked within last 24 hours
	if info, err := os.Stat(checkFile); err == nil {
		if time.Since(info.ModTime()) < 24*time.Hour {
			return
		}
	}

	// Touch the file regardless of outcome so we don't retry on failure
	os.MkdirAll(filepath.Dir(checkFile), 0755)
	os.WriteFile(checkFile, []byte(time.Now().Format(time.RFC3339)), 0644)

	latest, err := getLatestVersion()
	if err != nil {
		return
	}

	if isNewer(latest, Version) {
		fmt.Fprintf(os.Stderr, "dbx %s available (current: %s). Run 'dbx update' to upgrade.\n", latest, Version)
	}
}
