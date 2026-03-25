package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/chrisfentiman/dbx/internal/config"
	"github.com/chrisfentiman/dbx/internal/guard"
	"github.com/spf13/cobra"
)

var checkInputFile string

var checkCmd = &cobra.Command{
	Use:   `check ["SQL"]`,
	Short: "Validate SQL formatting and runnability against the guard",
	Long: `Checks a SQL query for formatting issues and whether the dbx guard would allow it.

  dbx check "SELECT * FROM catalog.schema.table"
  dbx check -f query.sql

Checks:
  - Guard validation (prohibited keywords, schema whitelist, parseability)
  - Formatting (keyword casing, clause separation)
  - Table references (fully qualified names)`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		sqlStr, err := resolveSQLArg(args)
		if err != nil {
			return err
		}

		issues := 0

		// 1. Guard validation
		cfg, cfgErr := config.Load()
		var g *guard.Guard
		if cfgErr == nil {
			g = guard.New(cfg.AllowedSchemas, cfg.DefaultCatalog)
		} else {
			g = guard.New(nil, "")
		}

		result := g.Validate(sqlStr)
		if result.Allowed {
			fmt.Println("✓ Guard: query is allowed")
			if result.Reason != "" {
				fmt.Printf("  %s\n", result.Reason)
			}
		} else {
			fmt.Printf("✗ Guard: %s\n", result.Reason)
			issues++
		}

		// 2. Table references
		if result.Allowed && len(result.Tables) > 0 {
			for _, t := range result.Tables {
				parts := strings.Split(t, ".")
				if len(parts) < 2 {
					fmt.Printf("⚠ Table %q is not fully qualified (use catalog.schema.table)\n", t)
					issues++
				}
			}
			if issues == 0 {
				fmt.Printf("✓ Tables: %s\n", strings.Join(result.Tables, ", "))
			}
		}

		// 3. Formatting check
		formatted := formatSQL(sqlStr)
		if strings.TrimSpace(formatted) != strings.TrimSpace(sqlStr) {
			fmt.Println("⚠ Formatting: query has style issues (run 'dbx fmt' to fix)")
			issues++
		} else {
			fmt.Println("✓ Formatting: clean")
		}

		// 4. Common issues
		upper := strings.ToUpper(sqlStr)
		if strings.Contains(upper, "SELECT *") {
			fmt.Println("⚠ Style: avoid SELECT * — list columns explicitly")
			issues++
		}
		if !strings.Contains(upper, "LIMIT") && strings.Contains(upper, "SELECT") {
			fmt.Println("⚠ Style: consider adding LIMIT for exploratory queries")
		}

		if issues > 0 {
			fmt.Printf("\n%d issue(s) found\n", issues)
			if !result.Allowed {
				return fmt.Errorf("query would be blocked")
			}
		}

		return nil
	},
}

func resolveSQLArg(args []string) (string, error) {
	// Check both fmt and check input file flags
	inputFile := fmtInputFile
	if checkInputFile != "" {
		inputFile = checkInputFile
	}

	switch {
	case inputFile != "" && len(args) > 0:
		return "", fmt.Errorf("provide SQL as an argument or with -f, not both")
	case inputFile != "":
		data, err := os.ReadFile(inputFile)
		if err != nil {
			return "", fmt.Errorf("reading file: %w", err)
		}
		return strings.TrimSpace(string(data)), nil
	case len(args) == 1:
		return args[0], nil
	default:
		return "", fmt.Errorf("provide SQL as an argument or with -f <file>")
	}
}

func init() {
	checkCmd.Flags().StringVarP(&checkInputFile, "file", "f", "", "read SQL from file")
	rootCmd.AddCommand(checkCmd)
}
