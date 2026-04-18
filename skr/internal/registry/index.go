package registry

import (
	"encoding/json"
	"os"
	"strings"
)

// Index represents the marketplace.json structure.
type Index struct {
	Packages []Package `json:"packages"`
}

// Package represents a single package entry in the index.
type Package struct {
	Name          string            `json:"name"`
	Type          string            `json:"type"`
	Version       string            `json:"version"`
	Description   string            `json:"description"`
	Tags          []string          `json:"tags"`
	Author        string            `json:"author"`
	License       string            `json:"license"`
	Path          string            `json:"path"`
	Files         []string          `json:"files"`
	ArtifactHash  string            `json:"artifact_hash,omitempty"`
	Notes         string            `json:"notes,omitempty"`
	InstallTarget string            `json:"install_target,omitempty"`
	Upstream      *UpstreamRef      `json:"upstream,omitempty"`
}

// UpstreamRef tracks external package sources.
type UpstreamRef struct {
	URL  string `json:"url"`
	Path string `json:"path"`
	Ref  string `json:"ref"`
}

// LoadIndex reads and parses a marketplace.json file.
func LoadIndex(path string) (*Index, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var idx Index
	if err := json.Unmarshal(data, &idx); err != nil {
		return nil, err
	}
	return &idx, nil
}

// FindByName returns a package by exact name match.
func (idx *Index) FindByName(name string) *Package {
	for i := range idx.Packages {
		if idx.Packages[i].Name == name {
			return &idx.Packages[i]
		}
	}
	return nil
}

// FindByNameVersion returns a package matching name and optional version.
func (idx *Index) FindByNameVersion(name, version string) *Package {
	for i := range idx.Packages {
		p := &idx.Packages[i]
		if p.Name == name {
			if version == "" || p.Version == version {
				return p
			}
		}
	}
	return nil
}

// Search returns packages matching a query string against name, description, and tags.
// Matching is case-insensitive substring.
func (idx *Index) Search(query string) []Package {
	if query == "" {
		return idx.Packages
	}
	q := strings.ToLower(query)
	var results []Package
	for _, p := range idx.Packages {
		if matchesQuery(p, q) {
			results = append(results, p)
		}
	}
	return results
}

// FilterByTag returns packages that have the given tag.
func (idx *Index) FilterByTag(tag string) []Package {
	t := strings.ToLower(tag)
	var results []Package
	for _, p := range idx.Packages {
		for _, pt := range p.Tags {
			if strings.ToLower(pt) == t {
				results = append(results, p)
				break
			}
		}
	}
	return results
}

func matchesQuery(p Package, query string) bool {
	if strings.Contains(strings.ToLower(p.Name), query) {
		return true
	}
	if strings.Contains(strings.ToLower(p.Description), query) {
		return true
	}
	for _, tag := range p.Tags {
		if strings.Contains(strings.ToLower(tag), query) {
			return true
		}
	}
	return false
}
