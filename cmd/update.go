package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

const (
	repoOwner = "chrisfentiman"
	repoName  = "dbx"
)

type githubRelease struct {
	TagName string `json:"tag_name"`
}

func getLatestVersion() (string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", repoOwner, repoName)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("GitHub API returned %d", resp.StatusCode)
	}

	var release githubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", err
	}

	return strings.TrimPrefix(release.TagName, "v"), nil
}

func isNewer(latest, current string) bool {
	if current == "dev" {
		return false
	}
	return latest != current
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update dbx to the latest version",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Print("Checking for updates... ")

		latest, err := getLatestVersion()
		if err != nil {
			return fmt.Errorf("checking latest version: %w", err)
		}

		if !isNewer(latest, Version) {
			fmt.Printf("already up to date (%s)\n", Version)
			return nil
		}

		fmt.Printf("updating %s → %s\n", Version, latest)

		// Determine binary target name
		goos := runtime.GOOS
		goarch := runtime.GOARCH
		target := fmt.Sprintf("dbx-%s-%s", goos, goarch)
		url := fmt.Sprintf("https://github.com/%s/%s/releases/latest/download/%s.tar.gz", repoOwner, repoName, target)

		// Download to temp dir
		fmt.Printf("Downloading %s...\n", target)
		resp, err := http.Get(url)
		if err != nil {
			return fmt.Errorf("downloading: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			return fmt.Errorf("download failed: HTTP %d", resp.StatusCode)
		}

		tmpFile, err := os.CreateTemp("", "dbx-update-*.tar.gz")
		if err != nil {
			return fmt.Errorf("creating temp file: %w", err)
		}
		defer os.Remove(tmpFile.Name())

		if _, err := io.Copy(tmpFile, resp.Body); err != nil {
			return fmt.Errorf("saving download: %w", err)
		}
		tmpFile.Close()

		// Extract
		tmpDir, err := os.MkdirTemp("", "dbx-update-")
		if err != nil {
			return fmt.Errorf("creating temp dir: %w", err)
		}
		defer os.RemoveAll(tmpDir)

		extractCmd := exec.Command("tar", "xzf", tmpFile.Name(), "-C", tmpDir)
		if err := extractCmd.Run(); err != nil {
			return fmt.Errorf("extracting: %w", err)
		}

		// Find current binary path
		currentBin, err := os.Executable()
		if err != nil {
			return fmt.Errorf("finding current binary: %w", err)
		}

		// Replace current binary
		newBin := fmt.Sprintf("%s/%s", tmpDir, target)
		if err := os.Rename(newBin, currentBin); err != nil {
			// Rename fails across filesystems, fall back to copy
			src, err := os.ReadFile(newBin)
			if err != nil {
				return fmt.Errorf("reading new binary: %w", err)
			}
			if err := os.WriteFile(currentBin, src, 0755); err != nil {
				return fmt.Errorf("writing new binary: %w", err)
			}
		}

		fmt.Printf("Updated to %s\n", latest)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
