package registry

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/philippeVV/skill-registry/skr/internal/config"
)

// RepoDir returns the path to the cached registry clone.
func RepoDir() string {
	return filepath.Join(config.SkrDir(), "cache", "repo")
}

// EnsureRepo ensures a fresh local clone of the registry exists.
// If the clone doesn't exist, it creates one. If it exists but is stale
// (older than CacheTTL), it pulls the latest. Returns the path to the
// repo directory.
func EnsureRepo(registryURL string, forceFresh bool) (string, error) {
	repoDir := RepoDir()
	repoURL := normalizeRepoURL(registryURL)

	if isRepoPresent(repoDir) {
		if forceFresh || isRepoStale(repoDir) {
			log.Debug().Msg("updating cached repo")
			if err := pullRepo(repoDir); err != nil {
				// Pull failed — nuke and re-clone
				log.Warn().Err(err).Msg("pull failed, re-cloning")
				os.RemoveAll(repoDir)
				return cloneRepo(repoURL, repoDir)
			}
			touchRepo(repoDir)
			return repoDir, nil
		}
		log.Debug().Msg("using cached repo")
		return repoDir, nil
	}

	return cloneRepo(repoURL, repoDir)
}

// FetchIndex reads marketplace.json from the local repo clone, refreshing
// if stale. Caches a copy of the index separately for offline use.
func FetchIndex(registryURL string) (*Index, error) {
	log.Debug().Str("registry", registryURL).Msg("fetching index")

	repoDir, err := EnsureRepo(registryURL, false)
	if err != nil {
		return nil, fmt.Errorf("ensuring repo: %w", err)
	}

	indexPath := filepath.Join(repoDir, "marketplace.json")
	data, err := os.ReadFile(indexPath)
	if err != nil {
		return nil, fmt.Errorf("reading marketplace.json: %w", err)
	}

	// Cache the index JSON separately for offline/fast reads
	if err := SaveCache(data); err != nil {
		log.Warn().Err(err).Msg("failed to cache index")
	}

	return LoadIndex(indexPath)
}

// FetchPackage returns the path to a specific package directory within the
// local repo clone. No additional cloning — just reads from the cached repo.
func FetchPackage(registryURL, pkgPath string) (string, func(), error) {
	log.Debug().Str("package", pkgPath).Msg("fetching package")

	repoDir, err := EnsureRepo(registryURL, false)
	if err != nil {
		return "", nil, fmt.Errorf("ensuring repo: %w", err)
	}

	pkgDir := filepath.Join(repoDir, pkgPath)
	if _, err := os.Stat(pkgDir); err != nil {
		return "", nil, fmt.Errorf("package not found at %s: %w", pkgPath, err)
	}

	// No cleanup needed — the repo stays cached
	noop := func() {}
	return pkgDir, noop, nil
}

// isRepoPresent checks if the cached repo exists and has a .git directory.
func isRepoPresent(repoDir string) bool {
	_, err := os.Stat(filepath.Join(repoDir, ".git"))
	return err == nil
}

// isRepoStale checks if the repo was last refreshed more than CacheTTL ago.
func isRepoStale(repoDir string) bool {
	markerPath := filepath.Join(repoDir, ".skr-fetched")
	info, err := os.Stat(markerPath)
	if err != nil {
		return true // no marker = stale
	}
	return time.Since(info.ModTime()) > CacheTTL
}

// touchRepo updates the fetch timestamp marker.
func touchRepo(repoDir string) {
	markerPath := filepath.Join(repoDir, ".skr-fetched")
	os.WriteFile(markerPath, []byte{}, 0o644)
}

// cloneRepo creates a fresh shallow clone.
func cloneRepo(repoURL, dest string) (string, error) {
	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return "", fmt.Errorf("creating cache dir: %w", err)
	}

	log.Debug().Str("url", repoURL).Str("dest", dest).Msg("cloning registry")
	cmd := exec.Command("git", "clone", "--depth", "1", repoURL, dest)
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("cloning registry: %w", err)
	}

	touchRepo(dest)
	return dest, nil
}

// pullRepo fetches and resets to the latest remote state.
func pullRepo(repoDir string) error {
	fetch := exec.Command("git", "-C", repoDir, "fetch", "--depth", "1", "origin", "main")
	fetch.Stderr = os.Stderr
	if err := fetch.Run(); err != nil {
		return fmt.Errorf("fetching: %w", err)
	}

	reset := exec.Command("git", "-C", repoDir, "reset", "--hard", "origin/main")
	reset.Stderr = os.Stderr
	if err := reset.Run(); err != nil {
		return fmt.Errorf("resetting: %w", err)
	}

	return nil
}

// normalizeRepoURL converts a GitHub URL to a git-cloneable URL.
func normalizeRepoURL(url string) string {
	if strings.HasPrefix(url, "git@") || strings.HasSuffix(url, ".git") {
		return url
	}
	if strings.HasPrefix(url, "https://github.com/") {
		return url + ".git"
	}
	return url
}
