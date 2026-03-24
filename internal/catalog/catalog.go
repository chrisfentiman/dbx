package catalog

import (
	"context"
	"fmt"

	"github.com/databricks/databricks-sdk-go"
	sdkcatalog "github.com/databricks/databricks-sdk-go/service/catalog"
)

// CatalogInfo is a simplified catalog representation.
type CatalogInfo struct {
	Name    string `json:"name"`
	Comment string `json:"comment,omitempty"`
	Owner   string `json:"owner,omitempty"`
}

// SchemaInfo is a simplified schema representation.
type SchemaInfo struct {
	Name    string `json:"name"`
	Comment string `json:"comment,omitempty"`
	Owner   string `json:"owner,omitempty"`
}

// TableInfo is a simplified table representation for list operations.
type TableInfo struct {
	Name      string `json:"name"`
	TableType string `json:"table_type"`
	Comment   string `json:"comment,omitempty"`
}

// TableDetail is a full table description with columns.
type TableDetail struct {
	Name      string       `json:"name"`
	FullName  string       `json:"full_name"`
	TableType string       `json:"table_type"`
	Comment   string       `json:"comment,omitempty"`
	Columns   []ColumnInfo `json:"columns"`
}

// ColumnInfo describes a single column.
type ColumnInfo struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Position int    `json:"position"`
	Nullable bool   `json:"nullable"`
	Comment  string `json:"comment,omitempty"`
}

// ListCatalogs returns all catalogs the user can access.
func ListCatalogs(ctx context.Context, w *databricks.WorkspaceClient) ([]CatalogInfo, error) {
	cats, err := w.Catalogs.ListAll(ctx, sdkcatalog.ListCatalogsRequest{})
	if err != nil {
		return nil, fmt.Errorf("listing catalogs: %w", err)
	}
	if len(cats) == 0 {
		return nil, fmt.Errorf("no catalogs found (check permissions)")
	}

	result := make([]CatalogInfo, 0, len(cats))
	for _, cat := range cats {
		result = append(result, CatalogInfo{
			Name:    cat.Name,
			Comment: cat.Comment,
			Owner:   cat.Owner,
		})
	}
	return result, nil
}

// ListSchemas returns all schemas in a catalog.
func ListSchemas(ctx context.Context, w *databricks.WorkspaceClient, catalogName string) ([]SchemaInfo, error) {
	schemas, err := w.Schemas.ListAll(ctx, sdkcatalog.ListSchemasRequest{
		CatalogName: catalogName,
	})
	if err != nil {
		return nil, fmt.Errorf("listing schemas in %s: %w", catalogName, err)
	}

	result := make([]SchemaInfo, 0, len(schemas))
	for _, s := range schemas {
		result = append(result, SchemaInfo{
			Name:    s.Name,
			Comment: s.Comment,
			Owner:   s.Owner,
		})
	}
	return result, nil
}

// ListTables returns all tables in a catalog.schema.
func ListTables(ctx context.Context, w *databricks.WorkspaceClient, catalogName, schemaName string) ([]TableInfo, error) {
	tables, err := w.Tables.ListAll(ctx, sdkcatalog.ListTablesRequest{
		CatalogName: catalogName,
		SchemaName:  schemaName,
	})
	if err != nil {
		return nil, fmt.Errorf("listing tables in %s.%s: %w", catalogName, schemaName, err)
	}

	result := make([]TableInfo, 0, len(tables))
	for _, t := range tables {
		result = append(result, TableInfo{
			Name:      t.Name,
			TableType: string(t.TableType),
			Comment:   t.Comment,
		})
	}
	return result, nil
}

// DescribeTable returns full table detail including columns.
func DescribeTable(ctx context.Context, w *databricks.WorkspaceClient, fullName string) (*TableDetail, error) {
	t, err := w.Tables.GetByFullName(ctx, fullName)
	if err != nil {
		return nil, fmt.Errorf("get table %s: %w", fullName, err)
	}

	detail := &TableDetail{
		Name:      t.Name,
		FullName:  t.FullName,
		TableType: string(t.TableType),
		Comment:   t.Comment,
	}

	for _, col := range t.Columns {
		detail.Columns = append(detail.Columns, ColumnInfo{
			Name:     col.Name,
			Type:     string(col.TypeName),
			Position: col.Position,
			Nullable: col.Nullable,
			Comment:  col.Comment,
		})
	}

	return detail, nil
}
