package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/philippeVV/skill-registry/skr/internal/config"
	"github.com/philippeVV/skill-registry/skr/internal/hash"
	"github.com/philippeVV/skill-registry/skr/internal/installer"
	"github.com/philippeVV/skill-registry/skr/internal/lockfile"
	"github.com/philippeVV/skill-registry/skr/internal/registry"
	"github.com/philippeVV/skill-registry/skr/internal/telemetry"
	"github.com/philippeVV/skill-registry/skr/internal/ui"
)

var tagFilter string

var installCmd = &cobra.Command{
	Use:   "install <name>[@version]",
	Short: "Install a package from the registry",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runInstall,
}

func init() {
	installCmd.Flags().StringVar(&tagFilter, "tag", "", "install all packages with this tag")
	rootCmd.AddCommand(installCmd)
}

func runInstall(cmd *cobra.Command, args []string) error {
	idx, err := getIndex()
	if err != nil {
		return err
	}

	lf, err := lockfile.Load()
	if err != nil {
		return fmt.Errorf("loading lockfile: %w", err)
	}

	isFirstInstall := !lf.HasAnyPackages()

	var toInstall []registry.Package

	if tagFilter != "" {
		toInstall = idx.FilterByTag(tagFilter)
		if len(toInstall) == 0 {
			fmt.Println(ui.Warning.Render(fmt.Sprintf("No packages found with tag %q", tagFilter)))
			return nil
		}
	} else if len(args) == 1 {
		name, version := parseNameVersion(args[0])
		pkg := idx.FindByNameVersion(name, version)
		if pkg == nil {
			return fmt.Errorf("package %q not found in registry", args[0])
		}
		toInstall = []registry.Package{*pkg}
	} else {
		return fmt.Errorf("specify a package name or --tag")
	}

	ctx := cmd.Context()
	for _, pkg := range toInstall {
		if err := installOne(ctx, cfg, &pkg, lf); err != nil {
			fmt.Println(ui.Error.Render(fmt.Sprintf("Failed to install %s: %v", pkg.Name, err)))
			continue
		}
	}

	// On first install, bootstrap system skills
	if isFirstInstall {
		if err := bootstrapSystemSkills(ctx, idx, lf); err != nil {
			log.Warn().Err(err).Msg("failed to bootstrap system skills")
		}
	}

	return lf.Save()
}

func installOne(ctx context.Context, cfg *config.Config, pkg *registry.Package, lf *lockfile.Lockfile) error {
	fmt.Printf("  Installing %s@%s (%s)...\n", pkg.Name, pkg.Version, pkg.Type)

	// Fetch the package
	srcDir, cleanup, err := registry.FetchPackage(cfg.Registry, pkg.Path)
	if err != nil {
		return err
	}
	defer cleanup()

	// Install to the correct location
	location, err := installer.Install(cfg, pkg, srcDir)
	if err != nil {
		return err
	}

	// Compute hash of installed artifact
	var h string
	switch pkg.Type {
	case "skill", "knowledge":
		h, err = hash.Dir(location)
	default:
		h, err = hash.File(location)
	}
	if err != nil {
		log.Warn().Err(err).Str("package", pkg.Name).Msg("failed to hash installed artifact")
	}

	lf.Add(pkg.Name, lockfile.Entry{
		Version:      pkg.Version,
		Type:         pkg.Type,
		Location:     location,
		Hash:         h,
		RegistryHash: pkg.ArtifactHash,
	})

	fmt.Println(ui.Success.Render(fmt.Sprintf("  Installed %s → %s", pkg.Name, location)))
	telemetry.EmitInstall(ctx, pkg.Name, pkg.Version, pkg.Type, cfg.Registry)
	return nil
}

func bootstrapSystemSkills(ctx context.Context, idx *registry.Index, lf *lockfile.Lockfile) error {
	systemSkills := idx.FilterByTag("registry-core")
	for _, pkg := range systemSkills {
		if lf.IsInstalled(pkg.Name) {
			continue
		}
		fmt.Printf("\n  Bootstrapping system skill: %s\n", pkg.Name)
		if err := installOne(ctx, cfg, &pkg, lf); err != nil {
			return err
		}
		// Mark as system
		entry, _ := lf.Get(pkg.Name)
		entry.System = true
		lf.Add(pkg.Name, entry)
	}
	return nil
}

func parseNameVersion(s string) (string, string) {
	if i := strings.LastIndex(s, "@"); i > 0 {
		return s[:i], s[i+1:]
	}
	return s, ""
}

// getIndex returns the registry index, using cache if fresh.
func getIndex() (*registry.Index, error) {
	if registry.CacheIsFresh() {
		idx, err := registry.LoadCachedIndex()
		if err == nil {
			log.Debug().Msg("using cached index")
			return idx, nil
		}
		log.Debug().Err(err).Msg("cache read failed, fetching fresh")
	}

	idx, err := registry.FetchIndex(cfg.Registry)
	if err != nil {
		// Try stale cache as fallback
		log.Warn().Err(err).Msg("fetch failed, trying stale cache")
		idx, cacheErr := registry.LoadCachedIndex()
		if cacheErr != nil {
			return nil, fmt.Errorf("fetch failed: %w (no cached index available)", err)
		}
		return idx, nil
	}
	return idx, nil
}
