package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/philippeVV/skill-registry/skr/internal/lockfile"
	"github.com/philippeVV/skill-registry/skr/internal/registry"
	"github.com/philippeVV/skill-registry/skr/internal/ui"
)

var updateCmd = &cobra.Command{
	Use:   "update [name]",
	Short: "Update installed packages to the latest version",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runUpdate,
}

func init() {
	rootCmd.AddCommand(updateCmd)
}

func runUpdate(cmd *cobra.Command, args []string) error {
	// Force refresh the repo and index
	if _, err := registry.EnsureRepo(cfg.Registry, true); err != nil {
		return fmt.Errorf("refreshing registry: %w", err)
	}
	idx, err := registry.FetchIndex(cfg.Registry)
	if err != nil {
		return fmt.Errorf("fetching index: %w", err)
	}

	lf, err := lockfile.Load()
	if err != nil {
		return fmt.Errorf("loading lockfile: %w", err)
	}

	var toUpdate []string
	if len(args) == 1 {
		toUpdate = []string{args[0]}
	} else {
		for name := range lf.Packages {
			toUpdate = append(toUpdate, name)
		}
	}

	updated := 0
	for _, name := range toUpdate {
		entry, ok := lf.Get(name)
		if !ok {
			fmt.Println(ui.Warning.Render(fmt.Sprintf("  %s is not installed, skipping", name)))
			continue
		}

		pkg := idx.FindByName(name)
		if pkg == nil {
			fmt.Println(ui.Warning.Render(fmt.Sprintf("  %s not found in registry, skipping", name)))
			continue
		}

		if pkg.Version == entry.Version {
			fmt.Printf("  %s is up to date (%s)\n", name, entry.Version)
			continue
		}

		fmt.Printf("  Updating %s: %s → %s\n", name, entry.Version, pkg.Version)
		if err := installOne(cfg, pkg, lf); err != nil {
			fmt.Println(ui.Error.Render(fmt.Sprintf("  Failed to update %s: %v", name, err)))
			continue
		}

		// Preserve system flag
		if entry.System {
			newEntry, _ := lf.Get(name)
			newEntry.System = true
			lf.Add(name, newEntry)
		}
		updated++
	}

	if err := lf.Save(); err != nil {
		return fmt.Errorf("saving lockfile: %w", err)
	}

	if updated == 0 {
		fmt.Println(ui.Muted.Render("\n  Everything is up to date."))
	} else {
		fmt.Println(ui.Success.Render(fmt.Sprintf("\n  Updated %d package(s).", updated)))
	}

	return nil
}
