package cmd

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// ScaffoldFS is set by main.go to the embedded filesystem.
var ScaffoldFS fs.FS

var forceOverwrite bool

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Scaffold the current directory with dbx skills, rules, and CLAUDE.md",
	Long: `Sets up the current directory for use with dbx and Claude Code.

Creates:
  .claude/skills/dbx-*   — Claude Code skills for data exploration
  .claude/rules/          — Code quality and style rules
  CLAUDE.md               — Project instructions for Claude Code
  .env.example            — Template for Databricks configuration

Behavior:
  - Skills (dbx-*) are always written (these are managed by dbx)
  - Rules are inserted only if they don't already exist
  - CLAUDE.md is skipped if it exists (use --force to overwrite)
  - .env.example is skipped if it exists
  - Never deletes existing files
  - Checks for uv (Python package manager) and offers to install it`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if ScaffoldFS == nil {
			return fmt.Errorf("scaffold files not available")
		}

		written := 0
		skipped := 0

		err := fs.WalkDir(ScaffoldFS, ".", func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if d.IsDir() {
				return os.MkdirAll(path, 0755)
			}

			action := classifyFile(path)

			switch action {
			case writeAlways:
				// dbx skills — always overwrite
			case writeIfMissing:
				// rules, .env.example — skip if exists
				if fileExists(path) {
					fmt.Printf("  skip  %s (already exists)\n", path)
					skipped++
					return nil
				}
			case writeIfForced:
				// CLAUDE.md — skip unless --force
				if fileExists(path) && !forceOverwrite {
					fmt.Printf("  skip  %s (use --force to overwrite)\n", path)
					skipped++
					return nil
				}
			}

			data, err := fs.ReadFile(ScaffoldFS, path)
			if err != nil {
				return fmt.Errorf("reading embedded %s: %w", path, err)
			}

			dir := filepath.Dir(path)
			if dir != "." {
				if err := os.MkdirAll(dir, 0755); err != nil {
					return fmt.Errorf("creating directory %s: %w", dir, err)
				}
			}

			if err := os.WriteFile(path, data, 0644); err != nil {
				return fmt.Errorf("writing %s: %w", path, err)
			}

			fmt.Printf("  write %s\n", path)
			written++
			return nil
		})

		if err != nil {
			return err
		}

		fmt.Printf("\nDone. %d files written, %d skipped.\n", written, skipped)

		// Check for Python data analysis dependencies
		checkPythonDeps()

		if !fileExists(".env") {
			fmt.Println("\nNext steps:")
			fmt.Println("  1. cp .env.example .env")
			fmt.Println("  2. Edit .env with your Databricks credentials")
			fmt.Println("  3. Run: dbx config setup")
		}

		return nil
	},
}

func checkPythonDeps() {
	// Check if uv is installed
	if _, err := exec.LookPath("uv"); err == nil {
		fmt.Println("\n✓ uv detected — Python scripts will use inline dependencies automatically")
		return
	}

	// Check if uvx is installed (sometimes uv installs as uvx)
	if _, err := exec.LookPath("uvx"); err == nil {
		fmt.Println("\n✓ uvx detected — Python scripts will use inline dependencies automatically")
		return
	}

	// uv not found — offer to install
	fmt.Println("\n⚠ uv not found")
	fmt.Println("  uv is a fast Python package manager that enables inline script dependencies.")
	fmt.Println("  Without it, Python data analysis (pandas, matplotlib, etc.) requires manual setup.")
	fmt.Print("\n  Install uv now? [Y/n] ")

	reader := bufio.NewReader(os.Stdin)
	answer, _ := reader.ReadString('\n')
	answer = strings.TrimSpace(strings.ToLower(answer))

	if answer == "" || answer == "y" || answer == "yes" {
		fmt.Println("  Installing uv...")
		installCmd := exec.Command("sh", "-c", "curl -LsSf https://astral.sh/uv/install.sh | sh")
		installCmd.Stdout = os.Stdout
		installCmd.Stderr = os.Stderr
		if err := installCmd.Run(); err != nil {
			fmt.Printf("  ✗ Failed to install uv: %v\n", err)
			fmt.Println("  Install manually: https://docs.astral.sh/uv/getting-started/installation/")
		} else {
			fmt.Println("  ✓ uv installed")
		}
	} else {
		fmt.Println("  Skipped. Install later: curl -LsSf https://astral.sh/uv/install.sh | sh")
	}
}

type fileAction int

const (
	writeAlways fileAction = iota
	writeIfMissing
	writeIfForced
)

func classifyFile(path string) fileAction {
	// dbx skills — always overwrite (managed by dbx)
	if matched, _ := filepath.Match(".claude/skills/dbx-*/*", path); matched {
		return writeAlways
	}

	// CLAUDE.md — only with --force
	if path == "CLAUDE.md" {
		return writeIfForced
	}

	// Everything else (rules, .env.example, non-dbx skills) — insert only
	return writeIfMissing
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func init() {
	upCmd.Flags().BoolVar(&forceOverwrite, "force", false, "overwrite CLAUDE.md if it already exists")
	rootCmd.AddCommand(upCmd)
}
