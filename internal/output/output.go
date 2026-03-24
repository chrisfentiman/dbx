package output

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"text/tabwriter"
)

// Result is the standard output envelope for all commands.
type Result struct {
	Command string      `json:"command"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Code    string      `json:"code,omitempty"`
}

// QueryResult specifically for SQL query results.
type QueryResult struct {
	Command  string       `json:"command"`
	SQL      string       `json:"sql,omitempty"`
	Columns  []ColumnMeta `json:"columns"`
	Rows     [][]string   `json:"rows"`
	RowCount int          `json:"row_count"`
}

// ColumnMeta describes a column in a query result set.
type ColumnMeta struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Nullable bool   `json:"nullable,omitempty"`
	Comment  string `json:"comment,omitempty"`
}

// WriteJSON writes any value as compact JSON to the writer.
func WriteJSON(w io.Writer, v interface{}) error {
	return json.NewEncoder(w).Encode(v)
}

// WriteTable writes a table with headers and rows using aligned columns.
func WriteTable(w io.Writer, headers []string, rows [][]string) error {
	tw := tabwriter.NewWriter(w, 0, 4, 2, ' ', 0)

	// Header row.
	for i, h := range headers {
		if i > 0 {
			fmt.Fprint(tw, "\t")
		}
		fmt.Fprint(tw, h)
	}
	fmt.Fprintln(tw)

	// Separator row.
	for i := range headers {
		if i > 0 {
			fmt.Fprint(tw, "\t")
		}
		fmt.Fprint(tw, "---")
	}
	fmt.Fprintln(tw)

	// Data rows.
	for _, row := range rows {
		for i, val := range row {
			if i > 0 {
				fmt.Fprint(tw, "\t")
			}
			fmt.Fprint(tw, val)
		}
		fmt.Fprintln(tw)
	}

	return tw.Flush()
}

// WriteCSV writes headers and rows as CSV to the writer.
func WriteCSV(w io.Writer, headers []string, rows [][]string) error {
	cw := csv.NewWriter(w)
	if err := cw.Write(headers); err != nil {
		return err
	}
	for _, row := range rows {
		if err := cw.Write(row); err != nil {
			return err
		}
	}
	cw.Flush()
	return cw.Error()
}

// WriteError writes an error message in the appropriate format.
func WriteError(w io.Writer, format string, code string, msg string) {
	switch format {
	case "json":
		WriteJSON(w, Result{Error: msg, Code: code})
	default:
		fmt.Fprintf(w, "Error [%s]: %s\n", code, msg)
	}
}

// Write dispatches to the appropriate formatter based on the format string.
// The data parameter is used for JSON output (should be a JSON-serializable struct).
// The headers and rows parameters are used for table and CSV output.
func Write(w io.Writer, format string, data interface{}, headers []string, rows [][]string) error {
	switch format {
	case "json":
		return WriteJSON(w, data)
	case "csv":
		return WriteCSV(w, headers, rows)
	case "table":
		return WriteTable(w, headers, rows)
	default:
		return fmt.Errorf("unknown format: %s", format)
	}
}
