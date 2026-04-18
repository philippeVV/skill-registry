package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/philippeVV/skill-registry/skr/internal/ui"
)

var searchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search for packages in the registry",
	Args:  cobra.ExactArgs(1),
	RunE:  runSearch,
}

func init() {
	rootCmd.AddCommand(searchCmd)
}

func runSearch(cmd *cobra.Command, args []string) error {
	query := args[0]

	idx, err := getIndex()
	if err != nil {
		return err
	}

	results := idx.Search(query)

	if len(results) == 0 {
		fmt.Println(ui.Warning.Render(fmt.Sprintf("  No packages matching %q", query)))
		return nil
	}

	headers := []string{"Name", "Type", "Description", "Tags"}
	var rows [][]string
	for _, p := range results {
		rows = append(rows, []string{
			p.Name,
			p.Type,
			truncate(p.Description, 50),
			strings.Join(p.Tags, ", "),
		})
	}

	fmt.Print(ui.PackageTable(headers, rows))
	return nil
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}
