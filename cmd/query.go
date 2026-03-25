package cmd

import (
	"fmt"
	"os"

	"github.com/chrisfentiman/dbx/internal/executor"
	"github.com/chrisfentiman/dbx/internal/guard"
	"github.com/chrisfentiman/dbx/internal/output"
	"github.com/spf13/cobra"
)

var (
	queryOutputFile string
	queryInputFile  string
)

var queryCmd = &cobra.Command{
	Use:   `query "SQL"`,
	Short: "Execute a read-only SQL query",
	Long: `Execute a read-only SQL query against Databricks.

Pass SQL directly as an argument:
  dbx query "SELECT * FROM schema.table LIMIT 10"

Or read SQL from a file:
  dbx query -f query.sql

Optionally write output to a file:
  dbx query "SELECT ..." --format csv -o results.csv`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		// Resolve SQL from args or file.
		var sqlStr string
		switch {
		case queryInputFile != "" && len(args) > 0:
			return fmt.Errorf("provide SQL as an argument or with -f, not both")
		case queryInputFile != "":
			data, err := os.ReadFile(queryInputFile)
			if err != nil {
				return fmt.Errorf("reading SQL file: %w", err)
			}
			sqlStr = string(data)
		case len(args) == 1:
			sqlStr = args[0]
		default:
			return fmt.Errorf("provide SQL as an argument or with -f <file>")
		}

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

		// Resolve output destination.
		w := os.Stdout
		if queryOutputFile != "" {
			f, err := os.Create(queryOutputFile)
			if err != nil {
				return fmt.Errorf("creating output file: %w", err)
			}
			defer f.Close()
			w = f
		}

		return output.Write(w, Format, data, headers, result.Rows)
	},
}

func init() {
	queryCmd.Flags().StringVarP(&queryOutputFile, "output", "o", "", "write output to file")
	queryCmd.Flags().StringVarP(&queryInputFile, "file", "f", "", "read SQL from file")
	rootCmd.AddCommand(queryCmd)
}
