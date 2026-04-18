package installer

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/philippeVV/skill-registry/skr/internal/config"
	"github.com/philippeVV/skill-registry/skr/internal/registry"
)

func installSkill(cfg *config.Config, pkg *registry.Package, srcDir string) (string, error) {
	destDir := filepath.Join(cfg.ClaudeConfigDir, "skills", pkg.Name)

	// Remove existing installation
	os.RemoveAll(destDir)

	if err := copyDir(srcDir, destDir); err != nil {
		return "", fmt.Errorf("installing skill %s: %w", pkg.Name, err)
	}
	return destDir, nil
}

func uninstallSkill(cfg *config.Config, pkg *registry.Package) error {
	destDir := filepath.Join(cfg.ClaudeConfigDir, "skills", pkg.Name)
	return os.RemoveAll(destDir)
}

// copyDir recursively copies a directory tree.
func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		// Skip metadata.json and README.md — registry-only files
		if rel == "metadata.json" || rel == "README.md" {
			return nil
		}

		destPath := filepath.Join(dst, rel)

		if info.IsDir() {
			return os.MkdirAll(destPath, 0o755)
		}

		return copyFile(path, destPath)
	})
}

func copyFile(src, dst string) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}

	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return out.Close()
}
