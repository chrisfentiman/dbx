package cmd

import (
	"fmt"
	"strings"

	"github.com/christopher-fentiman/dbx/internal/client"
	"github.com/christopher-fentiman/dbx/internal/config"
	"github.com/spf13/cobra"
)

// maskToken shows the first 4 and last 4 characters, masking the rest.
func maskToken(token string) string {
	if len(token) <= 8 {
		return "****"
	}
	return token[:4] + strings.Repeat("*", len(token)-8) + token[len(token)-4:]
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Show current configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("loading config: %w", err)
		}

		fmt.Printf("Host:            %s\n", cfg.Host)
		fmt.Printf("Token:           %s\n", maskToken(cfg.Token))
		fmt.Printf("Warehouse ID:    %s\n", cfg.WarehouseID)
		fmt.Printf("Default Catalog: %s\n", cfg.DefaultCatalog)
		fmt.Printf("Allowed Schemas: %s\n", strings.Join(cfg.AllowedSchemas, ", "))
		fmt.Printf("Default Format:  %s\n", cfg.DefaultFormat)
		fmt.Printf("Query Timeout:   %s\n", cfg.QueryTimeout)
		fmt.Printf("Max Rows:        %d\n", cfg.MaxRows)
		return nil
	},
}

var configValidateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate configuration and test connectivity",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("loading config: %w", err)
		}

		if err := cfg.Validate(); err != nil {
			return err
		}

		c, err := client.New(cfg)
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		me, err := c.W.CurrentUser.Me(cmd.Context())
		if err != nil {
			return fmt.Errorf("connectivity check failed: %w", err)
		}

		fmt.Printf("Connection successful! Logged in as: %s\n", me.UserName)
		return nil
	},
}

func init() {
	configCmd.AddCommand(configValidateCmd)
	rootCmd.AddCommand(configCmd)
}
