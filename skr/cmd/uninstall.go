package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/philippeVV/skill-registry/skr/internal/installer"
	"github.com/philippeVV/skill-registry/skr/internal/lockfile"
	"github.com/philippeVV/skill-registry/skr/internal/registry"
	"github.com/philippeVV/skill-registry/skr/internal/ui"
)

var forceUninstall bool

var uninstallCmd = &cobra.Command{
	Use:   "uninstall <name>",
	Short: "Uninstall a package",
	Args:  cobra.ExactArgs(1),
	RunE:  runUninstall,
}

func init() {
	uninstallCmd.Flags().BoolVar(&forceUninstall, "force", false, "force uninstall of system packages")
	rootCmd.AddCommand(uninstallCmd)
}

func runUninstall(cmd *cobra.Command, args []string) error {
	name := args[0]

	lf, err := lockfile.Load()
	if err != nil {
		return fmt.Errorf("loading lockfile: %w", err)
	}

	entry, ok := lf.Get(name)
	if !ok {
		return fmt.Errorf("package %q is not installed", name)
	}

	if entry.System && !forceUninstall {
		return fmt.Errorf("package %q is a system package — use --force to uninstall", name)
	}

	pkg := &registry.Package{
		Name: name,
		Type: entry.Type,
	}

	if err := installer.Uninstall(cfg, pkg); err != nil {
		return fmt.Errorf("uninstalling %s: %w", name, err)
	}

	lf.Remove(name)
	if err := lf.Save(); err != nil {
		return fmt.Errorf("saving lockfile: %w", err)
	}

	fmt.Println(ui.Success.Render(fmt.Sprintf("  Uninstalled %s", name)))
	return nil
}
