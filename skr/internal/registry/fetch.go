package registry

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
)

// FetchIndex clones the registry repo (shallow) and reads marketplace.json.
// Returns the parsed index and caches the raw JSON locally.
func FetchIndex(registryURL string) (*Index, error) {
	log.Debug().Str("registry", registryURL).Msg("fetching index")

	tmpDir, err := os.MkdirTemp("", "skr-fetch-*")
	if err != nil {
		return nil, fmt.Errorf("creating temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	repoURL := normalizeRepoURL(registryURL)
	if err := shallowClone(repoURL, tmpDir); err != nil {
		return nil, fmt.Errorf("cloning registry: %w", err)
	}

	indexPath := filepath.Join(tmpDir, "marketplace.json")
	data, err := os.ReadFile(indexPath)
	if err != nil {
		return nil, fmt.Errorf("reading marketplace.json: %w", err)
	}

	// Cache the index
	if err := SaveCache(data); err != nil {
		log.Warn().Err(err).Msg("failed to cache index")
	}

	return LoadIndex(indexPath)
}

// FetchPackage clones the registry and returns the path to a specific package directory.
// The caller is responsible for calling the returned cleanup function.
func FetchPackage(registryURL, pkgPath string) (string, func(), error) {
	log.Debug().Str("package", pkgPath).Msg("fetching package")

	tmpDir, err := os.MkdirTemp("", "skr-pkg-*")
	if err != nil {
		return "", nil, fmt.Errorf("creating temp dir: %w", err)
	}
	cleanup := func() { os.RemoveAll(tmpDir) }

	repoURL := normalizeRepoURL(registryURL)
	if err := shallowClone(repoURL, tmpDir); err != nil {
		cleanup()
		return "", nil, fmt.Errorf("cloning registry: %w", err)
	}

	extracted := filepath.Join(tmpDir, pkgPath)
	if _, err := os.Stat(extracted); err != nil {
		cleanup()
		return "", nil, fmt.Errorf("package not found at %s: %w", pkgPath, err)
	}

	return extracted, cleanup, nil
}

// shallowClone performs a simple shallow clone of the repo.
func shallowClone(repoURL, dest string) error {
	cmd := exec.Command("git", "clone", "--depth", "1", repoURL, dest)
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// normalizeRepoURL converts a GitHub URL to a git-cloneable URL.
func normalizeRepoURL(url string) string {
	// Already a git URL
	if strings.HasPrefix(url, "git@") || strings.HasSuffix(url, ".git") {
		return url
	}
	// GitHub HTTPS URL without .git suffix
	if strings.HasPrefix(url, "https://github.com/") {
		return url + ".git"
	}
	return url
}
