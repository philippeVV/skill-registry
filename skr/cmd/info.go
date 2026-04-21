package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/philippeVV/skill-registry/skr/internal/registry"
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

	pairs = append(pairs, []string{"Install target", resolveInstallTarget(pkg)})

	fmt.Println()
	fmt.Println(ui.Title.Render(fmt.Sprintf("  %s@%s", pkg.Name, pkg.Version)))
	fmt.Println()
	fmt.Print(ui.InfoBlock(pairs))
	fmt.Println()

	// Trust signals
	fmt.Println(ui.Bold.Render("  Trust Signals"))
	if cfg.OTELEndpoint != "" {
		fmt.Println(ui.Muted.Render("    Stats unavailable — coming soon"))
	} else {
		fmt.Println(ui.Muted.Render("    Stats unavailable — set otel_endpoint in config to see live counts."))
	}
	fmt.Println()

	// README (fall back to artifact file if no README.md)
	srcDir, cleanup, err := registry.FetchPackage(cfg.Registry, pkg.Path)
	if err == nil {
		defer cleanup()
		readmeData, readErr := os.ReadFile(filepath.Join(srcDir, "README.md"))
		if readErr != nil {
			artifactFile := map[string]string{
				"skill": "SKILL.md", "rule": "RULE.md", "knowledge": "KNOWLEDGE.md",
			}[pkg.Type]
			if artifactFile != "" {
				readmeData, readErr = os.ReadFile(filepath.Join(srcDir, artifactFile))
			}
		}
		if readErr == nil {
			fmt.Println(ui.Bold.Render("  README"))
			fmt.Println()
			fmt.Println(string(readmeData))
		}
	}

	return nil
}

// resolveInstallTarget returns a human-readable install target path for display.
func resolveInstallTarget(pkg *registry.Package) string {
	if pkg.InstallTarget != "" {
		return pkg.InstallTarget
	}
	switch pkg.Type {
	case "skill":
		return fmt.Sprintf("~/.claude/skills/%s/", pkg.Name)
	case "rule":
		return fmt.Sprintf("~/.claude/rules/%s.md", pkg.Name)
	case "knowledge":
		return fmt.Sprintf("~/.claude/knowledge/%s/", pkg.Name)
	default:
		return pkg.Type
	}
}
