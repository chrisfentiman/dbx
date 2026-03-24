package executor

import (
	"context"
	"fmt"
	"time"

	"github.com/databricks/databricks-sdk-go"
	"github.com/databricks/databricks-sdk-go/service/sql"
)

// Result holds the output of a SQL query execution.
type Result struct {
	Columns  []ColumnMeta `json:"columns"`
	Rows     [][]string   `json:"rows"`
	RowCount int          `json:"row_count"`
}

// ColumnMeta describes a result column.
type ColumnMeta struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// Execute runs a SQL statement against Databricks and returns the result.
func Execute(ctx context.Context, w *databricks.WorkspaceClient, warehouseID string, sqlStr string, timeout time.Duration) (*Result, error) {
	waitTimeout := fmt.Sprintf("%ds", int(timeout.Seconds()))

	resp, err := w.StatementExecution.ExecuteAndWait(ctx, sql.ExecuteStatementRequest{
		Statement:   sqlStr,
		WarehouseId: warehouseID,
		Format:      sql.FormatJsonArray,
		Disposition: sql.DispositionInline,
		WaitTimeout: waitTimeout,
	})
	if err != nil {
		return nil, fmt.Errorf("statement execution: %w", err)
	}

	// Safety check: ExecuteAndWait handles failures internally, but guard
	// against unexpected states in case the SDK behavior changes.
	if resp.Status != nil && resp.Status.State == sql.StatementStateFailed {
		msg := "query failed"
		if resp.Status.Error != nil {
			msg = resp.Status.Error.Message
		}
		return nil, fmt.Errorf("query failed: %s", msg)
	}

	result := &Result{}

	// Extract column metadata from the manifest schema.
	if resp.Manifest != nil && resp.Manifest.Schema != nil {
		for _, col := range resp.Manifest.Schema.Columns {
			result.Columns = append(result.Columns, ColumnMeta{
				Name: col.Name,
				Type: string(col.TypeName),
			})
		}
	}

	// Extract row data.
	if resp.Result != nil && resp.Result.DataArray != nil {
		result.Rows = resp.Result.DataArray
		result.RowCount = len(resp.Result.DataArray)
	}

	return result, nil
}
