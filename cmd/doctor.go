package cmd

import (
	"fmt"
	"strings"

	"github.com/chrisfentiman/dbx/internal/client"
	"github.com/chrisfentiman/dbx/internal/config"
	"github.com/spf13/cobra"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check configuration, connectivity, and permissions",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("dbx doctor")
		fmt.Println(strings.Repeat("─", 40))

		// 1. Config
		fmt.Print("Config file... ")
		cfg, err := config.Load()
		if err != nil {
			fmt.Printf("✗ %v\n", err)
			return nil
		}
		fmt.Println("✓ loaded")

		// 2. Required fields
		fmt.Print("Required fields... ")
		if err := cfg.Validate(); err != nil {
			fmt.Printf("✗ %v\n", err)
			return nil
		}
		fmt.Println("✓ set")

		// 3. Connectivity
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

		// 4. Summary
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

func init() {
	rootCmd.AddCommand(doctorCmd)
}
