package cmd

import (
	"os"

	"github.com/christopher-fentiman/dbx/internal/catalog"
	"github.com/christopher-fentiman/dbx/internal/output"
	"github.com/spf13/cobra"
)

var catalogsCmd = &cobra.Command{
	Use:   "catalogs",
	Short: "List accessible catalogs",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		cats, err := catalog.ListCatalogs(ctx, Client.W)
		if err != nil {
			return err
		}

		// Filter to only the default catalog if one is configured.
		if Cfg.DefaultCatalog != "" {
			var filtered []catalog.CatalogInfo
			for _, c := range cats {
				if c.Name == Cfg.DefaultCatalog {
					filtered = append(filtered, c)
				}
			}
			cats = filtered
		}

		headers := []string{"Name", "Comment", "Owner"}
		var rows [][]string
		for _, c := range cats {
			rows = append(rows, []string{c.Name, c.Comment, c.Owner})
		}

		data := struct {
			Command  string                `json:"command"`
			Catalogs []catalog.CatalogInfo `json:"catalogs"`
		}{Command: "catalogs", Catalogs: cats}

		return output.Write(os.Stdout, Format, data, headers, rows)
	},
}

func init() {
	rootCmd.AddCommand(catalogsCmd)
}
