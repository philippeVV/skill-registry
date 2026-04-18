package lockfile

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/philippeVV/skill-registry/skr/internal/config"
)

// Lockfile tracks installed packages.
type Lockfile struct {
	Registry string             `json:"registry"`
	Packages map[string]Entry   `json:"packages"`
}

// Entry represents a single installed package.
type Entry struct {
	Version      string `json:"version"`
	Type         string `json:"type"`
	Location     string `json:"location"`
	Hash         string `json:"hash"`
	RegistryHash string `json:"registry_hash"`
	InstalledAt  string `json:"installed_at"`
	System       bool   `json:"system"`
}

// Load reads the lockfile from disk. Returns an empty lockfile if it doesn't exist.
func Load() (*Lockfile, error) {
	path := config.LockPath()
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Lockfile{
				Packages: make(map[string]Entry),
			}, nil
		}
		return nil, err
	}
	var lf Lockfile
	if err := json.Unmarshal(data, &lf); err != nil {
		return nil, err
	}
	if lf.Packages == nil {
		lf.Packages = make(map[string]Entry)
	}
	return &lf, nil
}

// Save writes the lockfile to disk.
func (lf *Lockfile) Save() error {
	path := config.LockPath()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(lf, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

// Add records an installed package.
func (lf *Lockfile) Add(name string, entry Entry) {
	if entry.InstalledAt == "" {
		entry.InstalledAt = time.Now().UTC().Format(time.RFC3339)
	}
	lf.Packages[name] = entry
}

// Remove deletes a package from the lockfile.
func (lf *Lockfile) Remove(name string) {
	delete(lf.Packages, name)
}

// Get returns a package entry and whether it exists.
func (lf *Lockfile) Get(name string) (Entry, bool) {
	e, ok := lf.Packages[name]
	return e, ok
}

// IsInstalled returns true if the package is in the lockfile.
func (lf *Lockfile) IsInstalled(name string) bool {
	_, ok := lf.Packages[name]
	return ok
}

// HasAnyPackages returns true if any packages are installed.
func (lf *Lockfile) HasAnyPackages() bool {
	return len(lf.Packages) > 0
}

// IsDrifted returns true if the installed hash differs from the registry hash.
func IsDrifted(entry Entry) bool {
	return entry.Hash != "" && entry.RegistryHash != "" && entry.Hash != entry.RegistryHash
}
