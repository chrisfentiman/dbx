package cmd

import (
	"fmt"
	"sort"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/chrisfentiman/dbx/internal/catalog"
	"github.com/chrisfentiman/dbx/internal/client"
	"github.com/chrisfentiman/dbx/internal/config"
	"github.com/spf13/cobra"
)

var configAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add configuration items",
}

var configAddSchemaCmd = &cobra.Command{
	Use:   "schema [schema_name]",
	Short: "Add a schema to the allowed list",
	Long: `Add a schema to the allowed_schemas list without re-running full setup.

If a schema name is provided, it is validated against Databricks before adding.
If no argument is given, shows an interactive selector of available schemas.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("loading config: %w", err)
		}
		if err := cfg.Validate(); err != nil {
			return fmt.Errorf("config incomplete — run 'dbx config setup' first: %w", err)
		}

		existingSet := make(map[string]bool, len(cfg.AllowedSchemas))
		for _, s := range cfg.AllowedSchemas {
			existingSet[s] = true
		}

		catalogName := cfg.DefaultCatalog
		if catalogName == "" {
			return fmt.Errorf("no default_catalog configured — run 'dbx config setup' first")
		}

		// Connect to validate schemas exist.
		c, err := client.New(cfg)
		if err != nil {
			return fmt.Errorf("client creation failed: %w", err)
		}
		if _, err := c.W.CurrentUser.Me(ctx); err != nil {
			return fmt.Errorf("connection failed: %w", err)
		}

		schemas, err := catalog.ListSchemas(ctx, c.W, catalogName)
		if err != nil {
			return fmt.Errorf("listing schemas: %w", err)
		}

		schemaOptions := make([]string, len(schemas))
		validSet := make(map[string]bool, len(schemas))
		for i, s := range schemas {
			schemaOptions[i] = s.Name
			validSet[s.Name] = true
		}

		// Direct mode: validate the named schema exists.
		if len(args) == 1 {
			name := args[0]
			if existingSet[name] {
				fmt.Printf("Schema '%s' is already in the allowed list.\n", name)
				return nil
			}
			if !validSet[name] {
				return fmt.Errorf("schema '%s' does not exist in catalog '%s'", name, catalogName)
			}
			cfg.AllowedSchemas = append(cfg.AllowedSchemas, name)
			sort.Strings(cfg.AllowedSchemas)

			if err := cfg.Save(); err != nil {
				return err
			}
			fmt.Printf("Added schema '%s'.\n", name)
			fmt.Printf("Allowed schemas: %s\n", strings.Join(cfg.AllowedSchemas, ", "))
			return nil
		}

		// Pre-check currently allowed schemas.
		var schemaDefaults []string
		for _, opt := range schemaOptions {
			if existingSet[opt] {
				schemaDefaults = append(schemaDefaults, opt)
			}
		}

		var selected []string
		prompt := &survey.MultiSelect{
			Message:  "Select schemas to allow (space to toggle, enter to confirm):",
			Options:  schemaOptions,
			Default:  schemaDefaults,
			PageSize: 20,
		}
		if err := survey.AskOne(prompt, &selected); err != nil {
			return fmt.Errorf("schema selection: %w", err)
		}

		// Merge: keep schemas not in the remote list (from other catalogs).
		remoteSet := make(map[string]bool, len(schemaOptions))
		for _, s := range schemaOptions {
			remoteSet[s] = true
		}
		mergedSet := make(map[string]bool)
		for _, s := range cfg.AllowedSchemas {
			if !remoteSet[s] {
				mergedSet[s] = true
			}
		}
		for _, s := range selected {
			mergedSet[s] = true
		}

		merged := make([]string, 0, len(mergedSet))
		for s := range mergedSet {
			merged = append(merged, s)
		}
		sort.Strings(merged)

		// Report changes.
		var added []string
		for _, s := range merged {
			if !existingSet[s] {
				added = append(added, s)
			}
		}
		var removed []string
		for _, s := range cfg.AllowedSchemas {
			if !mergedSet[s] {
				removed = append(removed, s)
			}
		}

		cfg.AllowedSchemas = merged
		if err := cfg.Save(); err != nil {
			return err
		}

		fmt.Println()
		if len(added) > 0 {
			fmt.Printf("Added:   %s\n", strings.Join(added, ", "))
		}
		if len(removed) > 0 {
			fmt.Printf("Removed: %s\n", strings.Join(removed, ", "))
		}
		if len(added) == 0 && len(removed) == 0 {
			fmt.Println("No changes.")
		}
		fmt.Printf("Allowed schemas: %s\n", strings.Join(cfg.AllowedSchemas, ", "))
		return nil
	},
}

func init() {
	configAddCmd.AddCommand(configAddSchemaCmd)
	configCmd.AddCommand(configAddCmd)
}
