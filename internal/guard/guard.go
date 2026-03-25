package guard

import (
	"fmt"
	"strings"

	"github.com/xwb1989/sqlparser"
)

var DeniedKeywords = []string{
	"INSERT", "UPDATE", "DELETE", "DROP", "CREATE", "ALTER",
	"TRUNCATE", "MERGE", "REPLACE", "GRANT", "REVOKE",
	"COPY", "EXEC", "EXECUTE",
}

type Guard struct {
	AllowedSchemas map[string]bool
	DefaultCatalog string
}

type ValidationResult struct {
	Allowed bool     `json:"allowed"`
	Reason  string   `json:"reason,omitempty"`
	Tables  []string `json:"tables,omitempty"`
}

func New(allowedSchemas []string, defaultCatalog string) *Guard {
	m := make(map[string]bool, len(allowedSchemas))
	for _, s := range allowedSchemas {
		m[strings.ToLower(strings.TrimSpace(s))] = true
	}
	return &Guard{AllowedSchemas: m, DefaultCatalog: defaultCatalog}
}

func (g *Guard) Validate(sql string) ValidationResult {
	// Layer 1: Keyword scan (strip comments for prefix detection)
	cleaned := stripSQLComments(sql)
	upper := strings.ToUpper(strings.TrimSpace(cleaned))
	for _, kw := range DeniedKeywords {
		if containsKeyword(upper, kw) {
			return ValidationResult{Allowed: false, Reason: fmt.Sprintf("prohibited keyword: %s", kw)}
		}
	}

	// Layer 2: AST parse
	stmt, err := sqlparser.Parse(sql)
	if err != nil {
		// Parser doesn't understand Databricks-specific syntax.
		// Allow only if the query starts with a known read-only command.
		if isReadOnlyPrefix(upper) {
			return ValidationResult{Allowed: true, Reason: "keyword-check only (Databricks-specific syntax)"}
		}
		return ValidationResult{Allowed: false, Reason: "unable to parse query; ensure it is a SELECT, SHOW, or DESCRIBE statement"}
	}

	// Verify statement type
	switch stmt.(type) {
	case *sqlparser.Select, *sqlparser.Union, *sqlparser.ParenSelect:
		// Read-only OK
	case *sqlparser.Show, *sqlparser.OtherRead:
		// Metadata queries OK
	default:
		return ValidationResult{Allowed: false, Reason: fmt.Sprintf("only SELECT/SHOW/DESCRIBE/EXPLAIN allowed, got: %T", stmt)}
	}

	// Extract table references
	tables := extractTableNames(stmt)

	// Schema whitelist check (skip if no schemas configured)
	if len(g.AllowedSchemas) > 0 {
		for _, table := range tables {
			if !g.isSchemaAllowed(table) {
				return ValidationResult{
					Allowed: false,
					Reason:  fmt.Sprintf("table %q references a schema not in the allowed list", table),
					Tables:  tables,
				}
			}
		}
	}

	return ValidationResult{Allowed: true, Tables: tables}
}

func extractTableNames(stmt sqlparser.Statement) []string {
	var tables []string
	_ = sqlparser.Walk(func(node sqlparser.SQLNode) (bool, error) {
		switch n := node.(type) {
		case sqlparser.TableName:
			qualifier := n.Qualifier.String()
			name := n.Name.String()
			if qualifier != "" {
				tables = append(tables, qualifier+"."+name)
			} else if name != "" {
				tables = append(tables, name)
			}
		}
		return true, nil
	}, stmt)
	return tables
}

func (g *Guard) isSchemaAllowed(tableRef string) bool {
	parts := strings.Split(strings.ToLower(tableRef), ".")
	switch len(parts) {
	case 3: // catalog.schema.table
		return g.AllowedSchemas[parts[0]+"."+parts[1]] || g.AllowedSchemas[parts[1]]
	case 2: // schema.table
		catSchema := g.DefaultCatalog + "." + parts[0]
		return g.AllowedSchemas[parts[0]] || g.AllowedSchemas[strings.ToLower(catSchema)]
	case 1: // bare table — allowed (uses default schema context)
		return true
	default:
		return false
	}
}

func stripSQLComments(sql string) string {
	// Strip block comments /* ... */
	for {
		start := strings.Index(sql, "/*")
		if start == -1 {
			break
		}
		end := strings.Index(sql[start+2:], "*/")
		if end == -1 {
			break
		}
		sql = sql[:start] + sql[start+2+end+2:]
	}

	// Strip line comments -- ...
	lines := strings.Split(sql, "\n")
	var result []string
	for _, line := range lines {
		if idx := strings.Index(line, "--"); idx >= 0 {
			line = line[:idx]
		}
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			result = append(result, line)
		}
	}
	return strings.Join(result, "\n")
}

var readOnlyPrefixes = []string{
	"SELECT", "WITH", "SHOW", "DESCRIBE", "DESC", "EXPLAIN",
}

func isReadOnlyPrefix(upper string) bool {
	for _, prefix := range readOnlyPrefixes {
		if strings.HasPrefix(upper, prefix) {
			// Ensure it's a word boundary (followed by space, newline, tab, or end)
			if len(upper) == len(prefix) {
				return true
			}
			next := upper[len(prefix)]
			if next == ' ' || next == '\n' || next == '\t' || next == '\r' || next == '(' {
				return true
			}
		}
	}
	return false
}

func containsKeyword(upper, kw string) bool {
	idx := 0
	for {
		pos := strings.Index(upper[idx:], kw)
		if pos == -1 {
			return false
		}
		absPos := idx + pos
		before := absPos == 0 || !isIdentChar(upper[absPos-1])
		after := absPos+len(kw) >= len(upper) || !isIdentChar(upper[absPos+len(kw)])
		if before && after {
			return true
		}
		idx = absPos + len(kw)
	}
}

func isIdentChar(b byte) bool {
	return (b >= 'A' && b <= 'Z') || (b >= 'a' && b <= 'z') || (b >= '0' && b <= '9') || b == '_'
}
