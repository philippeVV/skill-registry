package lockfile

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTempLockfile(t *testing.T) (string, func()) {
	t.Helper()
	dir := t.TempDir()
	// Override the lock path via env
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", dir)
	// Create the config dir
	os.MkdirAll(filepath.Join(dir, ".config", "skr"), 0o755)
	return dir, func() { os.Setenv("HOME", oldHome) }
}

func TestLoadEmpty(t *testing.T) {
	_, cleanup := setupTempLockfile(t)
	defer cleanup()

	lf, err := Load()
	require.NoError(t, err)
	assert.Empty(t, lf.Packages)
	assert.False(t, lf.HasAnyPackages())
}

func TestAddAndGet(t *testing.T) {
	_, cleanup := setupTempLockfile(t)
	defer cleanup()

	lf, err := Load()
	require.NoError(t, err)

	lf.Add("test-pkg", Entry{
		Version:      "0.0.1",
		Type:         "skill",
		Location:     "/tmp/test",
		Hash:         "sha256:aaa",
		RegistryHash: "sha256:aaa",
	})

	assert.True(t, lf.HasAnyPackages())
	assert.True(t, lf.IsInstalled("test-pkg"))

	entry, ok := lf.Get("test-pkg")
	assert.True(t, ok)
	assert.Equal(t, "0.0.1", entry.Version)
	assert.NotEmpty(t, entry.InstalledAt)
}

func TestSaveAndReload(t *testing.T) {
	_, cleanup := setupTempLockfile(t)
	defer cleanup()

	lf, err := Load()
	require.NoError(t, err)

	lf.Registry = "https://example.com/registry"
	lf.Add("test-pkg", Entry{
		Version:      "1.0.0",
		Type:         "rule",
		Location:     "/tmp/rules/test.md",
		Hash:         "sha256:bbb",
		RegistryHash: "sha256:bbb",
	})

	require.NoError(t, lf.Save())

	// Reload
	lf2, err := Load()
	require.NoError(t, err)
	assert.Equal(t, "https://example.com/registry", lf2.Registry)

	entry, ok := lf2.Get("test-pkg")
	assert.True(t, ok)
	assert.Equal(t, "1.0.0", entry.Version)
	assert.Equal(t, "rule", entry.Type)
}

func TestRemove(t *testing.T) {
	_, cleanup := setupTempLockfile(t)
	defer cleanup()

	lf, err := Load()
	require.NoError(t, err)

	lf.Add("pkg-a", Entry{Version: "1.0.0", Type: "skill"})
	lf.Add("pkg-b", Entry{Version: "2.0.0", Type: "rule"})

	lf.Remove("pkg-a")
	assert.False(t, lf.IsInstalled("pkg-a"))
	assert.True(t, lf.IsInstalled("pkg-b"))
}

func TestDriftDetection(t *testing.T) {
	entry := Entry{
		Hash:         "sha256:aaa",
		RegistryHash: "sha256:aaa",
	}
	assert.False(t, IsDrifted(entry))

	entry.Hash = "sha256:bbb"
	assert.True(t, IsDrifted(entry))
}

func TestDriftDetectionEmptyHash(t *testing.T) {
	entry := Entry{
		Hash:         "",
		RegistryHash: "sha256:aaa",
	}
	assert.False(t, IsDrifted(entry), "empty hash should not count as drift")
}
