package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/christopher-fentiman/dbx/internal/catalog"
	"github.com/christopher-fentiman/dbx/internal/output"
	"github.com/spf13/cobra"
)

var tablesCmd = &cobra.Command{
	Use:   "tables <catalog.schema>",
	Short: "List tables in a schema",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		parts := strings.SplitN(args[0], ".", 2)
		if len(parts) != 2 {
			return fmt.Errorf("expected format: catalog.schema (got %q)", args[0])
		}
		catalogName, schemaName := parts[0], parts[1]

		tables, err := catalog.ListTables(ctx, Client.W, catalogName, schemaName)
		if err != nil {
			return err
		}

		headers := []string{"Name", "Type", "Comment"}
		var rows [][]string
		for _, t := range tables {
			rows = append(rows, []string{t.Name, t.TableType, t.Comment})
		}

		data := struct {
			Command string              `json:"command"`
			Catalog string              `json:"catalog"`
			Schema  string              `json:"schema"`
			Tables  []catalog.TableInfo `json:"tables"`
		}{Command: "tables", Catalog: catalogName, Schema: schemaName, Tables: tables}

		return output.Write(os.Stdout, Format, data, headers, rows)
	},
}

func init() {
	rootCmd.AddCommand(tablesCmd)
}
