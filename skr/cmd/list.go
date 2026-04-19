package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/philippeVV/skill-registry/skr/internal/hash"
	"github.com/philippeVV/skill-registry/skr/internal/lockfile"
	"github.com/philippeVV/skill-registry/skr/internal/ui"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List installed packages",
	Args:  cobra.NoArgs,
	RunE:  runList,
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func runList(cmd *cobra.Command, args []string) error {
	lf, err := lockfile.Load()
	if err != nil {
		return fmt.Errorf("loading lockfile: %w", err)
	}

	if len(lf.Packages) == 0 {
		fmt.Println(ui.Muted.Render("  No packages installed."))
		return nil
	}

	headers := []string{"Name", "Version", "Type", "Status"}
	var rows [][]string

	for name, entry := range lf.Packages {
		status := ui.Success.Render("ok")

		// Check for drift
		var currentHash string
		var hashErr error
		switch entry.Type {
		case "skill", "knowledge":
			currentHash, hashErr = hash.Dir(entry.Location)
		default:
			currentHash, hashErr = hash.File(entry.Location)
		}

		if hashErr != nil {
			status = ui.Error.Render("missing")
		} else if currentHash != entry.Hash {
			status = ui.BadgeModified
		}

		if entry.System {
			name = name + " " + ui.BadgeSystem
		}

		rows = append(rows, []string{name, entry.Version, entry.Type, status})
	}

	fmt.Print(ui.PackageTable(headers, rows))
	return nil
}
