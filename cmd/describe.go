package cmd

import (
	"os"
	"strconv"

	"github.com/chrisfentiman/dbx/internal/catalog"
	"github.com/chrisfentiman/dbx/internal/output"
	"github.com/spf13/cobra"
)

var describeCmd = &cobra.Command{
	Use:   "describe <catalog.schema.table>",
	Short: "Describe a table's columns and types",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		fullName := args[0]

		detail, err := catalog.DescribeTable(ctx, Client.W, fullName)
		if err != nil {
			return err
		}

		headers := []string{"Position", "Name", "Type", "Nullable", "Comment"}
		var rows [][]string
		for _, col := range detail.Columns {
			nullable := "NO"
			if col.Nullable {
				nullable = "YES"
			}
			rows = append(rows, []string{
				strconv.Itoa(col.Position),
				col.Name,
				col.Type,
				nullable,
				col.Comment,
			})
		}

		data := struct {
			Command string               `json:"command"`
			Table   string               `json:"table"`
			Type    string               `json:"table_type"`
			Comment string               `json:"comment,omitempty"`
			Columns []catalog.ColumnInfo `json:"columns"`
		}{
			Command: "describe",
			Table:   detail.FullName,
			Type:    detail.TableType,
			Comment: detail.Comment,
			Columns: detail.Columns,
		}

		return output.Write(os.Stdout, Format, data, headers, rows)
	},
}

func init() {
	rootCmd.AddCommand(describeCmd)
}
