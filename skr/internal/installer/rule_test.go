package installer

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/philippeVV/skill-registry/skr/internal/registry"
)

func TestInstallRule(t *testing.T) {
	cfg := testConfig(t)

	srcDir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(srcDir, "RULE.md"), []byte("# Test Rule\n"), 0o644))

	pkg := &registry.Package{Name: "test-rule", Type: "rule"}
	location, err := Install(cfg, pkg, srcDir)
	require.NoError(t, err)

	expectedPath := filepath.Join(cfg.ClaudeConfigDir, "rules", "test-rule.md")
	assert.Equal(t, expectedPath, location)

	content, err := os.ReadFile(expectedPath)
	require.NoError(t, err)
	assert.Equal(t, "# Test Rule\n", string(content))
}

func TestUninstallRule(t *testing.T) {
	cfg := testConfig(t)

	srcDir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(srcDir, "RULE.md"), []byte("# Rule\n"), 0o644))
	pkg := &registry.Package{Name: "remove-rule", Type: "rule"}

	_, err := Install(cfg, pkg, srcDir)
	require.NoError(t, err)

	rulePath := filepath.Join(cfg.ClaudeConfigDir, "rules", "remove-rule.md")
	assert.FileExists(t, rulePath)

	require.NoError(t, Uninstall(cfg, pkg))
	assert.NoFileExists(t, rulePath)
}
