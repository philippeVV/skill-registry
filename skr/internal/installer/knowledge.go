package installer

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/philippeVV/skill-registry/skr/internal/config"
	"github.com/philippeVV/skill-registry/skr/internal/registry"
)

// KnowledgeIndex is the local index of installed knowledge packages.
type KnowledgeIndex struct {
	Packages []KnowledgeEntry `json:"packages"`
}

// KnowledgeEntry is a single entry in the knowledge index.
type KnowledgeEntry struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
	Path        string   `json:"path"`
}

func knowledgeDir(cfg *config.Config) string {
	return filepath.Join(cfg.ClaudeConfigDir, "knowledge")
}

func knowledgeIndexPath(cfg *config.Config) string {
	return filepath.Join(knowledgeDir(cfg), "index.json")
}

func installKnowledge(cfg *config.Config, pkg *registry.Package, srcDir string) (string, error) {
	destDir := filepath.Join(knowledgeDir(cfg), pkg.Name)

	// Remove existing installation
	os.RemoveAll(destDir)

	if err := os.MkdirAll(destDir, 0o755); err != nil {
		return "", fmt.Errorf("creating knowledge dir: %w", err)
	}

	srcPath := filepath.Join(srcDir, "KNOWLEDGE.md")
	destPath := filepath.Join(destDir, "KNOWLEDGE.md")

	if err := copyFile(srcPath, destPath); err != nil {
		return "", fmt.Errorf("installing knowledge %s: %w", pkg.Name, err)
	}

	// Update the knowledge index
	if err := addToKnowledgeIndex(cfg, pkg); err != nil {
		return "", fmt.Errorf("updating knowledge index: %w", err)
	}

	return destDir, nil
}

func uninstallKnowledge(cfg *config.Config, pkg *registry.Package) error {
	destDir := filepath.Join(knowledgeDir(cfg), pkg.Name)
	if err := os.RemoveAll(destDir); err != nil {
		return err
	}
	return removeFromKnowledgeIndex(cfg, pkg.Name)
}

func loadKnowledgeIndex(cfg *config.Config) (*KnowledgeIndex, error) {
	path := knowledgeIndexPath(cfg)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &KnowledgeIndex{}, nil
		}
		return nil, err
	}
	var idx KnowledgeIndex
	if err := json.Unmarshal(data, &idx); err != nil {
		return nil, err
	}
	return &idx, nil
}

func saveKnowledgeIndex(cfg *config.Config, idx *KnowledgeIndex) error {
	path := knowledgeIndexPath(cfg)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(idx, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

func addToKnowledgeIndex(cfg *config.Config, pkg *registry.Package) error {
	idx, err := loadKnowledgeIndex(cfg)
	if err != nil {
		return err
	}

	// Remove existing entry for this package
	filtered := make([]KnowledgeEntry, 0, len(idx.Packages))
	for _, e := range idx.Packages {
		if e.Name != pkg.Name {
			filtered = append(filtered, e)
		}
	}

	filtered = append(filtered, KnowledgeEntry{
		Name:        pkg.Name,
		Description: pkg.Description,
		Tags:        pkg.Tags,
		Path:        pkg.Name + "/KNOWLEDGE.md",
	})

	idx.Packages = filtered
	return saveKnowledgeIndex(cfg, idx)
}

func removeFromKnowledgeIndex(cfg *config.Config, name string) error {
	idx, err := loadKnowledgeIndex(cfg)
	if err != nil {
		return err
	}

	filtered := make([]KnowledgeEntry, 0, len(idx.Packages))
	for _, e := range idx.Packages {
		if e.Name != name {
			filtered = append(filtered, e)
		}
	}

	idx.Packages = filtered
	return saveKnowledgeIndex(cfg, idx)
}
