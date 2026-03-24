package client

import (
	"github.com/christopher-fentiman/dbx/internal/config"
	"github.com/databricks/databricks-sdk-go"
)

// Client wraps the Databricks SDK for SQL warehouse operations.
type Client struct {
	W   *databricks.WorkspaceClient
	Cfg *config.Config
}

// New creates a new Databricks client from the given configuration.
func New(cfg *config.Config) (*Client, error) {
	w, err := databricks.NewWorkspaceClient(&databricks.Config{
		Host:  cfg.Host,
		Token: cfg.Token,
	})
	if err != nil {
		return nil, err
	}
	return &Client{
		W:   w,
		Cfg: cfg,
	}, nil
}
