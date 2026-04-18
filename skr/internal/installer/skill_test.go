package installer

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/philippeVV/skill-registry/skr/internal/config"
	"github.com/philippeVV/skill-registry/skr/internal/registry"
)

func testConfig(t *testing.T) *config.Config {
	t.Helper()
	dir := t.TempDir()
	return &config.Config{
		ClaudeConfigDir: dir,
	}
}

func TestInstallSkill(t *testing.T) {
	cfg := testConfig(t)

	// Create a source package directory
	srcDir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(srcDir, "SKILL.md"), []byte("# Test Skill\n"), 0o644))

	pkg := &registry.Package{
		Name: "test-skill",
		Type: "skill",
	}

	location, err := Install(cfg, pkg, srcDir)
	require.NoError(t, err)

	expectedDir := filepath.Join(cfg.ClaudeConfigDir, "skills", "test-skill")
	assert.Equal(t, expectedDir, location)

	// Verify SKILL.md was copied
	content, err := os.ReadFile(filepath.Join(expectedDir, "SKILL.md"))
	require.NoError(t, err)
	assert.Equal(t, "# Test Skill\n", string(content))
}

func TestInstallSkillWithSubdirs(t *testing.T) {
	cfg := testConfig(t)

	srcDir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(srcDir, "SKILL.md"), []byte("# Skill\n"), 0o644))
	require.NoError(t, os.MkdirAll(filepath.Join(srcDir, "scripts"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(srcDir, "scripts", "helper.sh"), []byte("#!/bin/bash\n"), 0o644))
	// metadata.json and README.md should be skipped
	require.NoError(t, os.WriteFile(filepath.Join(srcDir, "metadata.json"), []byte("{}"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(srcDir, "README.md"), []byte("# Readme"), 0o644))

	pkg := &registry.Package{Name: "multi-skill", Type: "skill"}
	location, err := Install(cfg, pkg, srcDir)
	require.NoError(t, err)

	// SKILL.md and scripts/ should exist
	assert.FileExists(t, filepath.Join(location, "SKILL.md"))
	assert.FileExists(t, filepath.Join(location, "scripts", "helper.sh"))

	// metadata.json and README.md should NOT be copied
	assert.NoFileExists(t, filepath.Join(location, "metadata.json"))
	assert.NoFileExists(t, filepath.Join(location, "README.md"))
}

func TestUninstallSkill(t *testing.T) {
	cfg := testConfig(t)

	// Install first
	srcDir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(srcDir, "SKILL.md"), []byte("# Test\n"), 0o644))
	pkg := &registry.Package{Name: "remove-me", Type: "skill"}
	_, err := Install(cfg, pkg, srcDir)
	require.NoError(t, err)

	// Verify it exists
	assert.DirExists(t, filepath.Join(cfg.ClaudeConfigDir, "skills", "remove-me"))

	// Uninstall
	require.NoError(t, Uninstall(cfg, pkg))
	assert.NoDirExists(t, filepath.Join(cfg.ClaudeConfigDir, "skills", "remove-me"))
}
