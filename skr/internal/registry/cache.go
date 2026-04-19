package registry

import (
	"os"
	"path/filepath"
	"time"

	"github.com/philippeVV/skill-registry/skr/internal/config"
)

const CacheTTL = 1 * time.Hour

// CacheIsFresh returns true if the cached index JSON exists and is less than CacheTTL old.
func CacheIsFresh() bool {
	info, err := os.Stat(config.CachePath())
	if err != nil {
		return false
	}
	return time.Since(info.ModTime()) < CacheTTL
}

// LoadCachedIndex reads the cached marketplace.json.
func LoadCachedIndex() (*Index, error) {
	return LoadIndex(config.CachePath())
}

// SaveCache writes index data to the cache location.
func SaveCache(data []byte) error {
	path := config.CachePath()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}
