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

var sqlKeywords = []string{
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

var sqlClauses = []string{
	"SELECT ", "FROM ", "WHERE ", "GROUP BY ", "ORDER BY ",
	"HAVING ", "LIMIT ", "INNER JOIN ", "LEFT JOIN ", "RIGHT JOIN ",
	"FULL JOIN ", "CROSS JOIN ", "UNION ALL", "UNION ", "WITH ",
}

func formatSQL(sql string) string {
	sql = strings.TrimSpace(sql)

	// Process line by line to preserve comments
	lines := strings.Split(sql, "\n")
	var result []string
	inBlockComment := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Track block comments
		if inBlockComment {
			if idx := strings.Index(line, "*/"); idx >= 0 {
				inBlockComment = false
				// Format the part after the block comment ends
				after := line[idx+2:]
				if strings.TrimSpace(after) != "" {
					after = formatSQLFragment(after)
				}
				result = append(result, line[:idx+2]+after)
			} else {
				result = append(result, line)
			}
			continue
		}

		// Full line comment — preserve as-is
		if strings.HasPrefix(trimmed, "--") {
			result = append(result, line)
			continue
		}

		// Block comment start
		if strings.HasPrefix(trimmed, "/*") {
			if idx := strings.Index(line, "*/"); idx >= 0 {
				// Single-line block comment
				after := line[idx+2:]
				if strings.TrimSpace(after) != "" {
					after = formatSQLFragment(after)
				}
				result = append(result, line[:idx+2]+after)
			} else {
				inBlockComment = true
				result = append(result, line)
			}
			continue
		}

		// Line with inline comment — split, format code part only
		if idx := strings.Index(line, "--"); idx >= 0 {
			code := line[:idx]
			comment := line[idx:]
			result = append(result, formatSQLFragment(code)+comment)
			continue
		}

		// Pure code line
		result = append(result, formatSQLFragment(line))
	}

	return strings.Join(result, "\n")
}

// formatSQLFragment formats a fragment of SQL (no comments)
func formatSQLFragment(sql string) string {
	for _, kw := range sqlKeywords {
		sql = replaceKeywordFmt(sql, kw)
	}
	return sql
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
