package cmd

import (
	"bufio"
	"encoding/json"
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
  - Files are only written if content has changed
  - Never deletes existing files
  - Checks for uv and pyright, offers to install if missing`,
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

			data, err := fs.ReadFile(ScaffoldFS, path)
			if err != nil {
				return fmt.Errorf("reading embedded %s: %w", path, err)
			}

			// Skip if file exists and content is identical
			if fileExists(path) && !contentChanged(path, data) {
				skipped++
				return nil
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

		// Check for Python LSP
		checkPythonLSP()

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

func checkPythonLSP() {
	// Check if pyright-langserver is installed
	if _, err := exec.LookPath("pyright-langserver"); err == nil {
		fmt.Println("✓ pyright detected — Python LSP diagnostics available")
		ensureLSPConfig()
		return
	}

	// Check if npm/node is available for installation
	if _, err := exec.LookPath("npm"); err != nil {
		fmt.Println("⊘ pyright not found (npm not available — install Node.js first)")
		return
	}

	fmt.Println("⚠ pyright not found")
	fmt.Println("  pyright provides Python type checking and diagnostics for Claude Code.")
	fmt.Print("\n  Install pyright now? [Y/n] ")

	reader := bufio.NewReader(os.Stdin)
	answer, _ := reader.ReadString('\n')
	answer = strings.TrimSpace(strings.ToLower(answer))

	if answer == "" || answer == "y" || answer == "yes" {
		fmt.Println("  Installing pyright...")
		installCmd := exec.Command("npm", "install", "-g", "pyright")
		installCmd.Stdout = os.Stdout
		installCmd.Stderr = os.Stderr
		if err := installCmd.Run(); err != nil {
			fmt.Printf("  ✗ Failed to install pyright: %v\n", err)
			fmt.Println("  Install manually: npm install -g pyright")
		} else {
			fmt.Println("  ✓ pyright installed")
			ensureLSPConfig()
		}
	} else {
		fmt.Println("  Skipped. Install later: npm install -g pyright")
	}
}

func ensureLSPConfig() {
	settingsPath := ".claude/settings.local.json"

	// Read existing settings or start fresh
	var settings map[string]interface{}
	if data, err := os.ReadFile(settingsPath); err == nil {
		if err := json.Unmarshal(data, &settings); err != nil {
			settings = make(map[string]interface{})
		}
	} else {
		settings = make(map[string]interface{})
	}

	// Check if LSP is already configured
	if _, exists := settings["lspServers"]; exists {
		return
	}

	// Add LSP config and ENABLE_LSP_TOOL env var
	settings["lspServers"] = map[string]interface{}{
		"python": map[string]interface{}{
			"command":   "pyright-langserver",
			"args":      []string{"--stdio"},
			"transport": "stdio",
			"extensionToLanguage": map[string]string{
				"py": "python",
			},
		},
	}

	// Add env var to enable LSP tool
	env, ok := settings["env"].(map[string]interface{})
	if !ok {
		env = make(map[string]interface{})
	}
	env["ENABLE_LSP_TOOL"] = "1"
	settings["env"] = env

	// Write back
	os.MkdirAll(".claude", 0755)
	data, _ := json.MarshalIndent(settings, "", "  ")
	os.WriteFile(settingsPath, append(data, '\n'), 0644)
	fmt.Println("  ✓ LSP config written to .claude/settings.local.json")
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func contentChanged(path string, newData []byte) bool {
	existing, err := os.ReadFile(path)
	if err != nil {
		return true // can't read = treat as changed
	}
	return string(existing) != string(newData)
}

func init() {
	rootCmd.AddCommand(upCmd)
}
