package cmd

import (
	"fmt"
	"os"

	"github.com/chrisfentiman/dbx/internal/executor"
	"github.com/chrisfentiman/dbx/internal/guard"
	"github.com/chrisfentiman/dbx/internal/output"
	"github.com/spf13/cobra"
)

var sampleLimit int

var sampleCmd = &cobra.Command{
	Use:   "sample <catalog.schema.table>",
	Short: "Sample rows from a table",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		tableName := args[0]
		sqlStr := fmt.Sprintf("SELECT * FROM %s LIMIT %d", tableName, sampleLimit)

		// Validate through SQL guard.
		g := guard.New(Cfg.AllowedSchemas, Cfg.DefaultCatalog)
		v := g.Validate(sqlStr)
		if !v.Allowed {
			return fmt.Errorf("query blocked: %s", v.Reason)
		}

		result, err := executor.Execute(ctx, Client.W, Cfg.WarehouseID, sqlStr, Timeout)
		if err != nil {
			return err
		}

		headers := make([]string, len(result.Columns))
		for i, col := range result.Columns {
			headers[i] = col.Name
		}

		data := struct {
			Command  string                `json:"command"`
			Table    string                `json:"table"`
			SQL      string                `json:"sql"`
			Columns  []executor.ColumnMeta `json:"columns"`
			Rows     [][]string            `json:"rows"`
			RowCount int                   `json:"row_count"`
		}{
			Command:  "sample",
			Table:    tableName,
			SQL:      sqlStr,
			Columns:  result.Columns,
			Rows:     result.Rows,
			RowCount: result.RowCount,
		}

		return output.Write(os.Stdout, Format, data, headers, result.Rows)
	},
}

func init() {
	sampleCmd.Flags().IntVar(&sampleLimit, "limit", 10, "number of rows to sample")
	rootCmd.AddCommand(sampleCmd)
}
