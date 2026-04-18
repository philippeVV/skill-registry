package installer

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/philippeVV/skill-registry/skr/internal/registry"
)

func TestInstallKnowledge(t *testing.T) {
	cfg := testConfig(t)

	srcDir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(srcDir, "KNOWLEDGE.md"), []byte("# Domain Knowledge\n"), 0o644))

	pkg := &registry.Package{
		Name:        "test-knowledge",
		Type:        "knowledge",
		Description: "Test knowledge for testing",
		Tags:        []string{"test", "example"},
	}

	location, err := Install(cfg, pkg, srcDir)
	require.NoError(t, err)

	expectedDir := filepath.Join(cfg.ClaudeConfigDir, "knowledge", "test-knowledge")
	assert.Equal(t, expectedDir, location)

	// Verify KNOWLEDGE.md was copied
	content, err := os.ReadFile(filepath.Join(expectedDir, "KNOWLEDGE.md"))
	require.NoError(t, err)
	assert.Equal(t, "# Domain Knowledge\n", string(content))

	// Verify index.json was updated
	indexPath := filepath.Join(cfg.ClaudeConfigDir, "knowledge", "index.json")
	assert.FileExists(t, indexPath)

	indexData, err := os.ReadFile(indexPath)
	require.NoError(t, err)

	var idx KnowledgeIndex
	require.NoError(t, json.Unmarshal(indexData, &idx))
	assert.Len(t, idx.Packages, 1)
	assert.Equal(t, "test-knowledge", idx.Packages[0].Name)
	assert.Equal(t, "Test knowledge for testing", idx.Packages[0].Description)
	assert.Equal(t, []string{"test", "example"}, idx.Packages[0].Tags)
}

func TestInstallKnowledgeUpdatesExistingIndex(t *testing.T) {
	cfg := testConfig(t)

	// Install first package
	srcDir1 := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(srcDir1, "KNOWLEDGE.md"), []byte("# First\n"), 0o644))
	pkg1 := &registry.Package{Name: "first", Type: "knowledge", Description: "First", Tags: []string{"a"}}
	_, err := Install(cfg, pkg1, srcDir1)
	require.NoError(t, err)

	// Install second package
	srcDir2 := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(srcDir2, "KNOWLEDGE.md"), []byte("# Second\n"), 0o644))
	pkg2 := &registry.Package{Name: "second", Type: "knowledge", Description: "Second", Tags: []string{"b"}}
	_, err = Install(cfg, pkg2, srcDir2)
	require.NoError(t, err)

	// Verify index has both entries
	indexData, err := os.ReadFile(filepath.Join(cfg.ClaudeConfigDir, "knowledge", "index.json"))
	require.NoError(t, err)

	var idx KnowledgeIndex
	require.NoError(t, json.Unmarshal(indexData, &idx))
	assert.Len(t, idx.Packages, 2)
}

func TestUninstallKnowledge(t *testing.T) {
	cfg := testConfig(t)

	srcDir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(srcDir, "KNOWLEDGE.md"), []byte("# Remove Me\n"), 0o644))
	pkg := &registry.Package{Name: "remove-knowledge", Type: "knowledge", Description: "test", Tags: []string{"x"}}

	_, err := Install(cfg, pkg, srcDir)
	require.NoError(t, err)

	// Verify it was installed
	knowledgeDir := filepath.Join(cfg.ClaudeConfigDir, "knowledge", "remove-knowledge")
	assert.DirExists(t, knowledgeDir)

	// Uninstall
	require.NoError(t, Uninstall(cfg, pkg))
	assert.NoDirExists(t, knowledgeDir)

	// Verify index is empty
	indexData, err := os.ReadFile(filepath.Join(cfg.ClaudeConfigDir, "knowledge", "index.json"))
	require.NoError(t, err)

	var idx KnowledgeIndex
	require.NoError(t, json.Unmarshal(indexData, &idx))
	assert.Empty(t, idx.Packages)
}
