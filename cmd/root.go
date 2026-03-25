package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/chrisfentiman/dbx/internal/client"
	"github.com/chrisfentiman/dbx/internal/config"
	"github.com/spf13/cobra"
)

var (
	// Flags accessible by subcommands.
	Format  string
	Catalog string
	Timeout time.Duration

	// Cfg holds the loaded configuration.
	Cfg *config.Config

	// Client holds the Databricks workspace client.
	Client *client.Client
)

var rootCmd = &cobra.Command{
	Use:   "dbx",
	Short: "Read-only Databricks SQL CLI for safe data exploration",
	Long: `dbx provides controlled, read-only SQL access to Databricks.

Setup:
  dbx up              Scaffold Claude Code skills in current directory
  dbx config setup    Interactive Databricks connection setup
  dbx doctor          Check config, connectivity, and binary integrity

Explore:
  dbx schemas <catalog>              List schemas
  dbx tables <catalog.schema>        List tables
  dbx describe <catalog.schema.tbl>  Show columns and types
  dbx sample <catalog.schema.tbl>    Preview rows
  dbx preview <catalog.schema.tbl>   Describe + sample

Query:
  dbx query "SELECT ..."             Execute SQL
  dbx query -f query.sql             Execute SQL from file
  dbx query "..." -o out.csv         Write results to file`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Skip config loading for commands that manage their own setup.
		switch cmd.Name() {
		case "version", "help", "config", "validate", "setup", "add", "schema", "up", "doctor", "update":
			return nil
		}

		// Background version check (once per day, non-blocking)
		go checkForUpdate()

		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("loading config: %w", err)
		}

		// Apply flag overrides when explicitly set.
		if cmd.Flags().Changed("catalog") {
			cfg.DefaultCatalog = Catalog
		}
		if cmd.Flags().Changed("format") {
			cfg.DefaultFormat = Format
		}
		if cmd.Flags().Changed("timeout") {
			cfg.QueryTimeout = Timeout
		}

		Cfg = cfg

		c, err := client.New(cfg)
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}
		Client = c

		return nil
	},
}

func init() {
	rootCmd.PersistentFlags().StringVar(&Format, "format", "table", "output format (json/csv/table)")
	rootCmd.PersistentFlags().StringVar(&Catalog, "catalog", "", "override default catalog")
	rootCmd.PersistentFlags().DurationVar(&Timeout, "timeout", 30*time.Second, "query timeout")
}

// Execute runs the root command.
func Execute() {
	rootCmd.SilenceUsage = true
	rootCmd.SilenceErrors = true

	if err := rootCmd.Execute(); err != nil {
		if Format == "json" {
			fmt.Fprintf(os.Stderr, "{\"error\": %q}\n", err.Error())
		} else {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}
		os.Exit(1)
	}
}
