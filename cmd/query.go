package cmd

import (
	"fmt"
	"os"

	"github.com/chrisfentiman/dbx/internal/executor"
	"github.com/chrisfentiman/dbx/internal/guard"
	"github.com/chrisfentiman/dbx/internal/output"
	"github.com/spf13/cobra"
)

var queryCmd = &cobra.Command{
	Use:   `query "SQL"`,
	Short: "Execute a read-only SQL query",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		sqlStr := args[0]

		// Validate through SQL guard.
		g := guard.New(Cfg.AllowedSchemas, Cfg.DefaultCatalog)
		v := g.Validate(sqlStr)
		if !v.Allowed {
			return fmt.Errorf("query blocked: %s", v.Reason)
		}

		// Execute the statement.
		result, err := executor.Execute(ctx, Client.W, Cfg.WarehouseID, sqlStr, Timeout)
		if err != nil {
			return err
		}

		// Build headers from column metadata.
		headers := make([]string, len(result.Columns))
		for i, col := range result.Columns {
			headers[i] = col.Name
		}

		data := struct {
			Command  string                `json:"command"`
			SQL      string                `json:"sql"`
			Columns  []executor.ColumnMeta `json:"columns"`
			Rows     [][]string            `json:"rows"`
			RowCount int                   `json:"row_count"`
		}{
			Command:  "query",
			SQL:      sqlStr,
			Columns:  result.Columns,
			Rows:     result.Rows,
			RowCount: result.RowCount,
		}

		return output.Write(os.Stdout, Format, data, headers, result.Rows)
	},
}

func init() {
	rootCmd.AddCommand(queryCmd)
}
