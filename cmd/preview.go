package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/christopher-fentiman/dbx/internal/catalog"
	"github.com/christopher-fentiman/dbx/internal/executor"
	"github.com/christopher-fentiman/dbx/internal/guard"
	"github.com/christopher-fentiman/dbx/internal/output"
	"github.com/spf13/cobra"
)

var previewCmd = &cobra.Command{
	Use:   "preview <catalog.schema.table>",
	Short: "Preview a table (describe + sample)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		fullName := args[0]

		// Step 1: Describe the table schema.
		detail, err := catalog.DescribeTable(ctx, Client.W, fullName)
		if err != nil {
			return err
		}

		// Step 2: Sample rows.
		sqlStr := fmt.Sprintf("SELECT * FROM %s LIMIT 5", fullName)
		g := guard.New(Cfg.AllowedSchemas, Cfg.DefaultCatalog)
		v := g.Validate(sqlStr)

		var sampleResult *executor.Result
		if v.Allowed {
			sampleResult, err = executor.Execute(ctx, Client.W, Cfg.WarehouseID, sqlStr, Timeout)
			if err != nil {
				// Non-fatal: still show schema even if sample fails.
				fmt.Fprintf(os.Stderr, "Warning: sample query failed: %v\n", err)
			}
		}

		if Format == "json" {
			data := struct {
				Command string              `json:"command"`
				Table   string              `json:"table"`
				Type    string              `json:"table_type"`
				Comment string              `json:"comment,omitempty"`
				Columns []catalog.ColumnInfo `json:"columns"`
				Sample  *executor.Result     `json:"sample,omitempty"`
			}{
				Command: "preview",
				Table:   detail.FullName,
				Type:    detail.TableType,
				Comment: detail.Comment,
				Columns: detail.Columns,
				Sample:  sampleResult,
			}
			return output.WriteJSON(os.Stdout, data)
		}

		// Table/CSV format: print schema then sample.
		fmt.Fprintf(os.Stdout, "=== Table: %s (%s) ===\n", detail.FullName, detail.TableType)
		if detail.Comment != "" {
			fmt.Fprintf(os.Stdout, "Comment: %s\n", detail.Comment)
		}
		fmt.Fprintln(os.Stdout, "\n--- Schema ---")

		schemaHeaders := []string{"#", "Name", "Type", "Nullable", "Comment"}
		var schemaRows [][]string
		for _, col := range detail.Columns {
			nullable := "NO"
			if col.Nullable {
				nullable = "YES"
			}
			schemaRows = append(schemaRows, []string{
				strconv.Itoa(col.Position), col.Name, col.Type, nullable, col.Comment,
			})
		}
		if err := output.WriteTable(os.Stdout, schemaHeaders, schemaRows); err != nil {
			return err
		}

		if sampleResult != nil && sampleResult.RowCount > 0 {
			fmt.Fprintln(os.Stdout, "\n--- Sample Data ---")
			sampleHeaders := make([]string, len(sampleResult.Columns))
			for i, col := range sampleResult.Columns {
				sampleHeaders[i] = col.Name
			}
			if err := output.WriteTable(os.Stdout, sampleHeaders, sampleResult.Rows); err != nil {
				return err
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(previewCmd)
}
