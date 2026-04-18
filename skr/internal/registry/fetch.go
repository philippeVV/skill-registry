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
// Returns the raw JSON bytes and caches them locally.
func FetchIndex(registryURL string) (*Index, error) {
	log.Debug().Str("registry", registryURL).Msg("fetching index")

	tmpDir, err := os.MkdirTemp("", "skr-fetch-*")
	if err != nil {
		return nil, fmt.Errorf("creating temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	repoURL := normalizeRepoURL(registryURL)
	cmd := exec.Command("git", "clone", "--depth", "1", "--filter=blob:none", "--sparse", repoURL, tmpDir)
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("cloning registry: %w", err)
	}

	// Sparse checkout just marketplace.json
	sparseCmd := exec.Command("git", "-C", tmpDir, "sparse-checkout", "set", "marketplace.json")
	if err := sparseCmd.Run(); err != nil {
		return nil, fmt.Errorf("sparse checkout: %w", err)
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

// FetchPackage clones the registry and extracts a specific package directory.
// Returns the path to the extracted package in a temp directory.
// The caller is responsible for cleaning up the returned directory's parent.
func FetchPackage(registryURL, pkgPath string) (string, func(), error) {
	log.Debug().Str("package", pkgPath).Msg("fetching package")

	tmpDir, err := os.MkdirTemp("", "skr-pkg-*")
	if err != nil {
		return "", nil, fmt.Errorf("creating temp dir: %w", err)
	}
	cleanup := func() { os.RemoveAll(tmpDir) }

	repoURL := normalizeRepoURL(registryURL)
	cmd := exec.Command("git", "clone", "--depth", "1", "--filter=blob:none", "--sparse", repoURL, tmpDir)
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		cleanup()
		return "", nil, fmt.Errorf("cloning registry: %w", err)
	}

	sparseCmd := exec.Command("git", "-C", tmpDir, "sparse-checkout", "set", pkgPath)
	if err := sparseCmd.Run(); err != nil {
		cleanup()
		return "", nil, fmt.Errorf("sparse checkout: %w", err)
	}

	extracted := filepath.Join(tmpDir, pkgPath)
	if _, err := os.Stat(extracted); err != nil {
		cleanup()
		return "", nil, fmt.Errorf("package not found at %s: %w", pkgPath, err)
	}

	return extracted, cleanup, nil
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
