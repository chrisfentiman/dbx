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
	Short: "Set up current directory for dbx + Claude Code",
	Long: `Scaffolds the current directory with Claude Code skills, rules, and config.

Writes .claude/skills/dbx-*, .claude/rules/*, CLAUDE.md, and .env.example.
Only writes files whose content has changed — safe to run repeatedly.
Checks for uv (Python) and pyright (LSP) and offers to install if missing.`,
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

		// Ensure dbx CLI permissions in settings
		ensurePermissions()

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
	settings := readSettings(settingsPath)

	if _, exists := settings["lspServers"]; exists {
		return
	}

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

	env, ok := settings["env"].(map[string]interface{})
	if !ok {
		env = make(map[string]interface{})
	}
	env["ENABLE_LSP_TOOL"] = "1"
	settings["env"] = env

	writeSettings(settingsPath, settings)
	fmt.Println("  ✓ LSP config written to .claude/settings.local.json")
}

func ensurePermissions() {
	settingsPath := ".claude/settings.local.json"

	settings := readSettings(settingsPath)

	// Allowed: data exploration and analysis commands
	requiredAllow := []string{
		"Bash(dbx query:*)",
		"Bash(dbx query *)",
		"Bash(dbx describe:*)",
		"Bash(dbx describe *)",
		"Bash(dbx sample:*)",
		"Bash(dbx sample *)",
		"Bash(dbx preview:*)",
		"Bash(dbx preview *)",
		"Bash(dbx tables:*)",
		"Bash(dbx tables *)",
		"Bash(dbx schemas:*)",
		"Bash(dbx schemas *)",
		"Bash(dbx catalogs:*)",
		"Bash(dbx catalogs *)",
		"Bash(dbx version:*)",
		"Bash(dbx doctor:*)",
		"Bash(python3:*)",
		"Bash(python3 *)",
		"Bash(uv run:*)",
		"Bash(uv run *)",
		"Bash(pip3 list:*)",
		"Bash(pip3 install:*)",
	}

	// Denied: configuration and system commands the LLM should not run
	requiredDeny := []string{
		"Bash(dbx config:*)",
		"Bash(dbx config *)",
		"Bash(dbx config setup:*)",
		"Bash(dbx config setup *)",
		"Bash(dbx up:*)",
		"Bash(dbx up *)",
		"Bash(dbx up)",
		"Bash(dbx update:*)",
		"Bash(dbx update *)",
		"Bash(dbx update)",
	}

	// Get existing permissions
	perms, ok := settings["permissions"].(map[string]interface{})
	if !ok {
		perms = make(map[string]interface{})
	}

	allow, ok := perms["allow"].([]interface{})
	if !ok {
		allow = []interface{}{}
	}

	// Build sets of existing permissions
	existingAllow := make(map[string]bool)
	for _, p := range allow {
		if s, ok := p.(string); ok {
			existingAllow[s] = true
		}
	}

	deny, _ := perms["deny"].([]interface{})
	existingDeny := make(map[string]bool)
	for _, p := range deny {
		if s, ok := p.(string); ok {
			existingDeny[s] = true
		}
	}

	// Add missing allow permissions
	changed := 0
	for _, perm := range requiredAllow {
		if !existingAllow[perm] {
			allow = append(allow, perm)
			changed++
		}
	}

	// Add missing deny permissions
	for _, perm := range requiredDeny {
		if !existingDeny[perm] {
			deny = append(deny, perm)
			changed++
		}
	}

	if changed == 0 {
		return
	}

	perms["allow"] = allow
	perms["deny"] = deny
	settings["permissions"] = perms
	writeSettings(settingsPath, settings)
	fmt.Println("  ✓ CLI permissions configured in .claude/settings.local.json")
}

func readSettings(path string) map[string]interface{} {
	var settings map[string]interface{}
	if data, err := os.ReadFile(path); err == nil {
		if err := json.Unmarshal(data, &settings); err != nil {
			settings = make(map[string]interface{})
		}
	} else {
		settings = make(map[string]interface{})
	}
	return settings
}

func writeSettings(path string, settings map[string]interface{}) {
	os.MkdirAll(".claude", 0755)
	data, _ := json.MarshalIndent(settings, "", "  ")
	os.WriteFile(path, append(data, '\n'), 0644)
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
