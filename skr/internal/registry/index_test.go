package registry

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testdataPath() string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(file), "..", "..", "testdata", "marketplace.json")
}

func TestLoadIndex(t *testing.T) {
	idx, err := LoadIndex(testdataPath())
	require.NoError(t, err)
	assert.NotEmpty(t, idx.Packages)
}

func TestFindByName(t *testing.T) {
	idx, err := LoadIndex(testdataPath())
	require.NoError(t, err)

	pkg := idx.FindByName("suggest-packages")
	require.NotNil(t, pkg)
	assert.Equal(t, "skill", pkg.Type)
	assert.Equal(t, "0.0.1", pkg.Version)
}

func TestFindByNameNotFound(t *testing.T) {
	idx, err := LoadIndex(testdataPath())
	require.NoError(t, err)

	pkg := idx.FindByName("nonexistent")
	assert.Nil(t, pkg)
}

func TestFindByNameVersion(t *testing.T) {
	idx, err := LoadIndex(testdataPath())
	require.NoError(t, err)

	pkg := idx.FindByNameVersion("suggest-packages", "0.0.1")
	require.NotNil(t, pkg)

	pkg = idx.FindByNameVersion("suggest-packages", "99.99.99")
	assert.Nil(t, pkg)
}

func TestSearchByName(t *testing.T) {
	idx, err := LoadIndex(testdataPath())
	require.NoError(t, err)

	results := idx.Search("suggest")
	assert.NotEmpty(t, results)
	assert.Equal(t, "suggest-packages", results[0].Name)
}

func TestSearchByDescription(t *testing.T) {
	idx, err := LoadIndex(testdataPath())
	require.NoError(t, err)

	results := idx.Search("code review")
	assert.NotEmpty(t, results)
	found := false
	for _, r := range results {
		if r.Name == "code-review-checklist" {
			found = true
		}
	}
	assert.True(t, found)
}

func TestSearchByTag(t *testing.T) {
	idx, err := LoadIndex(testdataPath())
	require.NoError(t, err)

	results := idx.Search("registry-core")
	assert.Len(t, results, 2) // suggest-packages and knowledge-retriever
}

func TestSearchEmpty(t *testing.T) {
	idx, err := LoadIndex(testdataPath())
	require.NoError(t, err)

	results := idx.Search("")
	assert.Equal(t, len(idx.Packages), len(results))
}

func TestSearchNoResults(t *testing.T) {
	idx, err := LoadIndex(testdataPath())
	require.NoError(t, err)

	results := idx.Search("zzzznonexistentzzzz")
	assert.Empty(t, results)
}

func TestFilterByTag(t *testing.T) {
	idx, err := LoadIndex(testdataPath())
	require.NoError(t, err)

	results := idx.FilterByTag("registry-core")
	assert.Len(t, results, 2)

	results = idx.FilterByTag("go")
	assert.Len(t, results, 1)
	assert.Equal(t, "go-conventions", results[0].Name)
}

func TestSearchCaseInsensitive(t *testing.T) {
	idx, err := LoadIndex(testdataPath())
	require.NoError(t, err)

	r1 := idx.Search("SUGGEST")
	r2 := idx.Search("suggest")
	assert.Equal(t, len(r1), len(r2))
}
