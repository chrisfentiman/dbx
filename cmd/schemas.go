package cmd

import (
	"fmt"
	"os"

	"github.com/christopher-fentiman/dbx/internal/catalog"
	"github.com/christopher-fentiman/dbx/internal/output"
	"github.com/spf13/cobra"
)

var schemasCmd = &cobra.Command{
	Use:   "schemas [catalog]",
	Short: "List schemas in a catalog",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		catalogName := Cfg.DefaultCatalog
		if len(args) > 0 {
			catalogName = args[0]
		}
		if catalogName == "" {
			return fmt.Errorf("catalog name required (provide as argument or set default_catalog in config)")
		}

		schemas, err := catalog.ListSchemas(ctx, Client.W, catalogName)
		if err != nil {
			return err
		}

		// Filter to allowed schemas if configured.
		if len(Cfg.AllowedSchemas) > 0 {
			allowed := make(map[string]bool, len(Cfg.AllowedSchemas))
			for _, s := range Cfg.AllowedSchemas {
				allowed[s] = true
			}
			var filtered []catalog.SchemaInfo
			for _, s := range schemas {
				if allowed[s.Name] {
					filtered = append(filtered, s)
				}
			}
			schemas = filtered
		}

		headers := []string{"Name", "Comment", "Owner"}
		var rows [][]string
		for _, s := range schemas {
			rows = append(rows, []string{s.Name, s.Comment, s.Owner})
		}

		data := struct {
			Command string              `json:"command"`
			Catalog string              `json:"catalog"`
			Schemas []catalog.SchemaInfo `json:"schemas"`
		}{Command: "schemas", Catalog: catalogName, Schemas: schemas}

		return output.Write(os.Stdout, Format, data, headers, rows)
	},
}

func init() {
	rootCmd.AddCommand(schemasCmd)
}
