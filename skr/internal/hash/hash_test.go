package hash

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileHash(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.txt")
	require.NoError(t, os.WriteFile(path, []byte("hello world\n"), 0o644))

	h, err := File(path)
	require.NoError(t, err)
	assert.Contains(t, h, "sha256:")
	assert.Len(t, h, 7+64) // "sha256:" + 64 hex chars
}

func TestFileHashDeterministic(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.txt")
	require.NoError(t, os.WriteFile(path, []byte("deterministic"), 0o644))

	h1, err := File(path)
	require.NoError(t, err)

	h2, err := File(path)
	require.NoError(t, err)

	assert.Equal(t, h1, h2)
}

func TestDirHash(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "a.txt"), []byte("aaa"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "b.txt"), []byte("bbb"), 0o644))

	h, err := Dir(dir)
	require.NoError(t, err)
	assert.Contains(t, h, "sha256:")
}

func TestDirHashChangesWithContent(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "a.txt"), []byte("original"), 0o644))

	h1, err := Dir(dir)
	require.NoError(t, err)

	require.NoError(t, os.WriteFile(filepath.Join(dir, "a.txt"), []byte("modified"), 0o644))

	h2, err := Dir(dir)
	require.NoError(t, err)

	assert.NotEqual(t, h1, h2)
}
