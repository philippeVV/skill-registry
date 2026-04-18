package installer

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/philippeVV/skill-registry/skr/internal/config"
	"github.com/philippeVV/skill-registry/skr/internal/registry"
)

func installRule(cfg *config.Config, pkg *registry.Package, srcDir string) (string, error) {
	rulesDir := filepath.Join(cfg.ClaudeConfigDir, "rules")
	if err := os.MkdirAll(rulesDir, 0o755); err != nil {
		return "", fmt.Errorf("creating rules dir: %w", err)
	}

	srcPath := filepath.Join(srcDir, "RULE.md")
	destPath := filepath.Join(rulesDir, pkg.Name+".md")

	if err := copyFile(srcPath, destPath); err != nil {
		return "", fmt.Errorf("installing rule %s: %w", pkg.Name, err)
	}
	return destPath, nil
}

func uninstallRule(cfg *config.Config, pkg *registry.Package) error {
	destPath := filepath.Join(cfg.ClaudeConfigDir, "rules", pkg.Name+".md")
	if err := os.Remove(destPath); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}
