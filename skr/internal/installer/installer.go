package installer

import (
	"fmt"

	"github.com/philippeVV/skill-registry/skr/internal/config"
	"github.com/philippeVV/skill-registry/skr/internal/registry"
)

// Install places a package artifact at the correct location based on type.
// srcDir is the extracted package directory from the registry.
// Returns the install location path.
func Install(cfg *config.Config, pkg *registry.Package, srcDir string) (string, error) {
	target := resolveTarget(pkg)
	switch target {
	case "skills":
		return installSkill(cfg, pkg, srcDir)
	case "rules":
		return installRule(cfg, pkg, srcDir)
	case "knowledge":
		return installKnowledge(cfg, pkg, srcDir)
	default:
		return "", fmt.Errorf("unknown install target %q for type %q", target, pkg.Type)
	}
}

// Uninstall removes a package from the install location.
func Uninstall(cfg *config.Config, pkg *registry.Package) error {
	target := resolveTarget(pkg)
	switch target {
	case "skills":
		return uninstallSkill(cfg, pkg)
	case "rules":
		return uninstallRule(cfg, pkg)
	case "knowledge":
		return uninstallKnowledge(cfg, pkg)
	default:
		return fmt.Errorf("unknown install target %q for type %q", target, pkg.Type)
	}
}

// resolveTarget determines the install target for a package.
// Priority: install_target override > type default.
func resolveTarget(pkg *registry.Package) string {
	if pkg.InstallTarget != "" {
		return pkg.InstallTarget
	}
	switch pkg.Type {
	case "skill":
		return "skills"
	case "rule":
		return "rules"
	case "knowledge":
		return "knowledge"
	default:
		return pkg.Type
	}
}
