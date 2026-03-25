package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/chrisfentiman/dbx/internal/client"
	"github.com/chrisfentiman/dbx/internal/config"
	"github.com/spf13/cobra"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Verify config, connectivity, version, and binary integrity",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("dbx doctor")
		fmt.Println(strings.Repeat("─", 40))

		// 1. Binary integrity
		fmt.Print("Binary integrity... ")
		if Version == "dev" {
			fmt.Println("⊘ dev build (skipped)")
		} else {
			if err := verifyBinaryIntegrity(); err != nil {
				fmt.Printf("✗ %v\n", err)
			} else {
				fmt.Println("✓ matches official release")
			}
		}

		// 2. Config
		fmt.Print("Config file... ")
		cfg, err := config.Load()
		if err != nil {
			fmt.Printf("✗ %v\n", err)
			return nil
		}
		fmt.Println("✓ loaded")

		// 3. Required fields
		fmt.Print("Required fields... ")
		if err := cfg.Validate(); err != nil {
			fmt.Printf("✗ %v\n", err)
			return nil
		}
		fmt.Println("✓ set")

		// 4. Connectivity
		fmt.Print("Connectivity... ")
		c, err := client.New(cfg)
		if err != nil {
			fmt.Printf("✗ %v\n", err)
			return nil
		}

		me, err := c.W.CurrentUser.Me(cmd.Context())
		if err != nil {
			fmt.Printf("✗ %v\n", err)
			return nil
		}
		fmt.Printf("✓ %s\n", me.UserName)

		// 5. Version check
		fmt.Print("Latest version... ")
		latest, err := getLatestVersion()
		if err != nil {
			fmt.Printf("✗ %v\n", err)
		} else if isNewer(latest, Version) {
			fmt.Printf("⚠ %s available (run 'dbx update')\n", latest)
		} else {
			fmt.Printf("✓ up to date (%s)\n", Version)
		}

		// 6. Summary
		fmt.Println(strings.Repeat("─", 40))
		fmt.Printf("Host:      %s\n", cfg.Host)
		fmt.Printf("Warehouse: %s\n", cfg.WarehouseID)
		fmt.Printf("Catalog:   %s\n", cfg.DefaultCatalog)
		fmt.Printf("Schemas:   %s\n", strings.Join(cfg.AllowedSchemas, ", "))
		fmt.Printf("Format:    %s\n", cfg.DefaultFormat)
		fmt.Printf("Timeout:   %s\n", cfg.QueryTimeout)
		fmt.Printf("Version:   %s\n", Version)

		return nil
	},
}

func verifyBinaryIntegrity() error {
	binPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("cannot find binary: %w", err)
	}

	binData, err := os.ReadFile(binPath)
	if err != nil {
		return fmt.Errorf("cannot read binary: %w", err)
	}

	localHash := sha256sum(binData)

	// Download official tarball and its checksum
	target := fmt.Sprintf("dbx-%s-%s", runtime.GOOS, runtime.GOARCH)
	baseURL := fmt.Sprintf("https://github.com/%s/%s/releases/download/v%s", repoOwner, repoName, Version)

	checksumData, err := downloadFile(fmt.Sprintf("%s/%s.tar.gz.sha256", baseURL, target))
	if err != nil {
		return fmt.Errorf("cannot fetch checksum: %w", err)
	}

	tarballData, err := downloadFile(fmt.Sprintf("%s/%s.tar.gz", baseURL, target))
	if err != nil {
		return fmt.Errorf("cannot fetch official binary: %w", err)
	}

	// Verify tarball against published checksum
	expectedTarHash := strings.Fields(string(checksumData))[0]
	if sha256sum(tarballData) != expectedTarHash {
		return fmt.Errorf("official release checksum mismatch (release may be corrupted)")
	}

	// Extract and compare binary
	tmpDir, err := os.MkdirTemp("", "dbx-verify-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)

	tmpTar, err := os.CreateTemp("", "dbx-verify-*.tar.gz")
	if err != nil {
		return err
	}
	defer os.Remove(tmpTar.Name())

	tmpTar.Write(tarballData)
	tmpTar.Close()

	if err := exec.Command("tar", "xzf", tmpTar.Name(), "-C", tmpDir).Run(); err != nil {
		return fmt.Errorf("extracting: %w", err)
	}

	officialBin, err := os.ReadFile(fmt.Sprintf("%s/%s", tmpDir, target))
	if err != nil {
		return fmt.Errorf("reading official binary: %w", err)
	}

	if localHash != sha256sum(officialBin) {
		return fmt.Errorf("MISMATCH — local binary differs from official release")
	}

	return nil
}

func init() {
	rootCmd.AddCommand(doctorCmd)
}
