package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var fmtFix bool

var fmtCmd = &cobra.Command{
	Use:   `fmt ["SQL"]`,
	Short: "Format a SQL query with consistent style",
	Long: `Formats SQL with consistent keyword casing and clause separation.

  dbx fmt "select * from catalog.schema.table where x = 1"
  dbx fmt -f query.sql
  dbx fmt -f query.sql --fix    # overwrites the file in place`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		sqlStr, err := resolveSQLArg(args)
		if err != nil {
			return err
		}

		formatted := formatSQL(sqlStr)

		if fmtFix && fmtInputFile != "" {
			if err := os.WriteFile(fmtInputFile, []byte(formatted+"\n"), 0644); err != nil {
				return fmt.Errorf("writing file: %w", err)
			}
			fmt.Printf("Formatted %s\n", fmtInputFile)
		} else {
			fmt.Println(formatted)
		}

		return nil
	},
}

var fmtInputFile string

func formatSQL(sql string) string {
	sql = strings.TrimSpace(sql)

	keywords := []string{
		"SELECT", "FROM", "WHERE", "AND", "OR", "NOT",
		"JOIN", "INNER JOIN", "LEFT JOIN", "RIGHT JOIN", "FULL JOIN", "CROSS JOIN",
		"LEFT OUTER JOIN", "RIGHT OUTER JOIN", "FULL OUTER JOIN",
		"ON", "AS", "IN", "IS", "NULL", "NOT NULL",
		"GROUP BY", "ORDER BY", "HAVING", "LIMIT", "OFFSET",
		"UNION", "UNION ALL", "INTERSECT", "EXCEPT",
		"CASE", "WHEN", "THEN", "ELSE", "END",
		"DISTINCT", "ALL", "EXISTS",
		"BETWEEN", "LIKE", "ILIKE",
		"ASC", "DESC", "NULLS FIRST", "NULLS LAST",
		"WITH", "RECURSIVE",
		"TRUE", "FALSE",
		"COUNT", "SUM", "AVG", "MIN", "MAX",
		"COALESCE", "CAST", "ROUND",
		"OVER", "PARTITION BY", "ROW_NUMBER", "RANK", "DENSE_RANK",
	}

	result := sql
	for _, kw := range keywords {
		result = replaceKeywordFmt(result, kw)
	}

	// Add newlines before major clauses
	clauses := []string{
		"SELECT ", "FROM ", "WHERE ", "GROUP BY ", "ORDER BY ",
		"HAVING ", "LIMIT ", "INNER JOIN ", "LEFT JOIN ", "RIGHT JOIN ",
		"FULL JOIN ", "CROSS JOIN ", "UNION ALL", "UNION ", "WITH ",
	}

	for _, clause := range clauses {
		idx := strings.Index(strings.ToUpper(result), clause)
		if idx > 0 {
			result = result[:idx] + "\n" + strings.TrimLeft(result[idx:], " ")
		}
	}

	return result
}

func replaceKeywordFmt(sql, keyword string) string {
	lower := strings.ToLower(keyword)
	upper := strings.ToUpper(keyword)
	result := []byte(sql)
	sqlLower := strings.ToLower(sql)

	for i := 0; i <= len(sqlLower)-len(lower); i++ {
		if sqlLower[i:i+len(lower)] == lower {
			before := i == 0 || !fmtIsIdentChar(result[i-1])
			after := i+len(lower) >= len(result) || !fmtIsIdentChar(result[i+len(lower)])
			if before && after {
				copy(result[i:i+len(upper)], upper)
			}
		}
	}
	return string(result)
}

func fmtIsIdentChar(b byte) bool {
	return (b >= 'A' && b <= 'Z') || (b >= 'a' && b <= 'z') || (b >= '0' && b <= '9') || b == '_'
}

func init() {
	fmtCmd.Flags().StringVarP(&fmtInputFile, "file", "f", "", "read SQL from file")
	fmtCmd.Flags().BoolVar(&fmtFix, "fix", false, "overwrite the file in place")
	rootCmd.AddCommand(fmtCmd)
}
