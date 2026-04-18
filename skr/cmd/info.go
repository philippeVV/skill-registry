package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/philippeVV/skill-registry/skr/internal/ui"
)

var infoCmd = &cobra.Command{
	Use:   "info <name>",
	Short: "Show detailed information about a package",
	Args:  cobra.ExactArgs(1),
	RunE:  runInfo,
}

func init() {
	rootCmd.AddCommand(infoCmd)
}

func runInfo(cmd *cobra.Command, args []string) error {
	name := args[0]

	idx, err := getIndex()
	if err != nil {
		return err
	}

	pkg := idx.FindByName(name)
	if pkg == nil {
		return fmt.Errorf("package %q not found in registry", name)
	}

	pairs := [][]string{
		{"Name", pkg.Name},
		{"Type", pkg.Type},
		{"Version", pkg.Version},
		{"Description", pkg.Description},
		{"Tags", strings.Join(pkg.Tags, ", ")},
		{"Author", pkg.Author},
		{"License", pkg.License},
	}

	if pkg.Notes != "" {
		pairs = append(pairs, []string{"Notes", pkg.Notes})
	}

	if pkg.Upstream != nil {
		pairs = append(pairs, []string{"Upstream", pkg.Upstream.URL})
	}

	fmt.Println()
	fmt.Println(ui.Title.Render(fmt.Sprintf("  %s@%s", pkg.Name, pkg.Version)))
	fmt.Println()
	fmt.Print(ui.InfoBlock(pairs))
	fmt.Println()

	return nil
}
