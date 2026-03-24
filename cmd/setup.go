package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/chrisfentiman/dbx/internal/catalog"
	"github.com/chrisfentiman/dbx/internal/client"
	"github.com/chrisfentiman/dbx/internal/config"
	"github.com/spf13/cobra"
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Interactive setup for dbx configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		reader := bufio.NewReader(os.Stdin)
		ctx := cmd.Context()

		// Load existing config as defaults (if any).
		existing, _ := config.Load()

		fmt.Println("=== dbx setup ===")
		fmt.Println("This will create or update ~/.dbx/config.yaml")
		fmt.Println()

		// Step 1: Connection credentials (plain prompts — these are short inputs).
		host := prompt(reader, "Databricks host URL", existing.Host)
		token := prompt(reader, "Personal access token", "")
		if token == "" && existing.Token != "" {
			fmt.Println("  (keeping existing token)")
			token = existing.Token
		}
		warehouseID := prompt(reader, "SQL warehouse ID", existing.WarehouseID)

		// Build a temporary config to connect.
		tmpCfg := &config.Config{
			Host:        host,
			Token:       token,
			WarehouseID: warehouseID,
		}
		if err := tmpCfg.Validate(); err != nil {
			return fmt.Errorf("invalid config: %w", err)
		}

		// Step 2: Connect and verify.
		fmt.Println()
		fmt.Print("Connecting to Databricks... ")
		c, err := client.New(tmpCfg)
		if err != nil {
			return fmt.Errorf("client creation failed: %w", err)
		}
		me, err := c.W.CurrentUser.Me(ctx)
		if err != nil {
			return fmt.Errorf("connection failed: %w", err)
		}
		fmt.Printf("OK (logged in as %s)\n", me.UserName)
		fmt.Println()

		// Step 3: Fetch catalogs and let user pick one interactively.
		fmt.Print("Fetching catalogs... ")
		catalogs, err := catalog.ListCatalogs(ctx, c.W)
		if err != nil {
			return fmt.Errorf("listing catalogs: %w", err)
		}
		fmt.Printf("found %d\n", len(catalogs))

		catalogOptions := make([]string, len(catalogs))
		for i, cat := range catalogs {
			catalogOptions[i] = cat.Name
		}

		catalogDefault := ""
		for _, opt := range catalogOptions {
			if opt == existing.DefaultCatalog {
				catalogDefault = opt
				break
			}
		}

		var selectedCatalog string
		catalogPrompt := &survey.Select{
			Message:  "Default catalog:",
			Options:  catalogOptions,
			PageSize: 15,
		}
		if catalogDefault != "" {
			catalogPrompt.Default = catalogDefault
		}
		if err := survey.AskOne(catalogPrompt, &selectedCatalog); err != nil {
			return fmt.Errorf("catalog selection: %w", err)
		}
		defaultCatalog := selectedCatalog

		// Step 4: Fetch schemas and let user multi-select interactively.
		fmt.Println()
		fmt.Printf("Fetching schemas in '%s'... ", defaultCatalog)
		schemas, err := catalog.ListSchemas(ctx, c.W, defaultCatalog)
		if err != nil {
			return fmt.Errorf("listing schemas: %w", err)
		}
		fmt.Printf("found %d\n", len(schemas))

		schemaOptions := make([]string, len(schemas))
		for i, s := range schemas {
			schemaOptions[i] = s.Name
		}

		// Pre-select previously allowed schemas.
		var schemaDefaults []string
		existingSet := make(map[string]bool)
		for _, s := range existing.AllowedSchemas {
			existingSet[s] = true
		}
		for _, opt := range schemaOptions {
			if existingSet[opt] {
				schemaDefaults = append(schemaDefaults, opt)
			}
		}

		var selectedSchemas []string
		schemaPrompt := &survey.MultiSelect{
			Message:  "Allowed schemas (type to filter, space to toggle, enter to confirm):",
			Options:  schemaOptions,
			Default:  schemaDefaults,
			PageSize: 20,
		}
		if err := survey.AskOne(schemaPrompt, &selectedSchemas); err != nil {
			return fmt.Errorf("schema selection: %w", err)
		}

		allowedSchemas := selectedSchemas

		// Step 5: Build final config.
		cfg := &config.Config{
			Host:           host,
			Token:          token,
			WarehouseID:    warehouseID,
			DefaultCatalog: defaultCatalog,
			AllowedSchemas: allowedSchemas,
			DefaultFormat:  "table",
			QueryTimeout:   existing.QueryTimeout,
			MaxRows:        existing.MaxRows,
		}
		if cfg.QueryTimeout == 0 {
			cfg.QueryTimeout = 30 * 1e9 // 30s
		}
		if cfg.MaxRows == 0 {
			cfg.MaxRows = 10000
		}

		// Step 6: Write config file.
		if err := cfg.Save(); err != nil {
			return err
		}

		// Summary.
		fmt.Println()
		fmt.Println("=== Setup Complete ===")
		fmt.Printf("Host:              %s\n", cfg.Host)
		fmt.Printf("Default catalog:   %s\n", defaultCatalog)
		fmt.Printf("Allowed schemas:   %s\n", strings.Join(allowedSchemas, ", "))
		fmt.Println()
		fmt.Println("Run './dbx-cli config validate' to verify.")

		return nil
	},
}

func prompt(reader *bufio.Reader, label string, defaultVal string) string {
	if defaultVal != "" {
		fmt.Printf("%s [%s]: ", label, defaultVal)
	} else {
		fmt.Printf("%s: ", label)
	}
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if input == "" {
		return defaultVal
	}
	return input
}

func init() {
	configCmd.AddCommand(setupCmd)
}
