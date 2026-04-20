package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const (
	DefaultRegistry = "https://github.com/philippeVV/skill-registry"
	ConfigDir       = ".config/skr"
	ConfigFile      = "config.json"
	CacheDir        = "cache"
	LockFile        = "skr.lock"
)

// Config holds the skr configuration.
type Config struct {
	Registry        string `json:"registry"`
	ClaudeConfigDir string `json:"claude_config_dir"`
	OTELEndpoint    string `json:"otel_endpoint"`
}

// DefaultClaudeConfigDir returns the default Claude config directory.
func DefaultClaudeConfigDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".claude")
}

// Load reads the config from disk, falling back to defaults.
// Environment variables override file values.
func Load() (*Config, error) {
	cfg := &Config{
		Registry:       DefaultRegistry,
		ClaudeConfigDir: DefaultClaudeConfigDir(),
	}

	path := ConfigPath()
	data, err := os.ReadFile(path)
	if err == nil {
		if err := json.Unmarshal(data, cfg); err != nil {
			return nil, err
		}
	}

	// Env var overrides
	if v := os.Getenv("SKR_REGISTRY"); v != "" {
		cfg.Registry = v
	}
	if v := os.Getenv("SKR_CLAUDE_CONFIG_DIR"); v != "" {
		cfg.ClaudeConfigDir = v
	}
	if v := os.Getenv("SKR_OTEL_ENDPOINT"); v != "" {
		cfg.OTELEndpoint = v
	}

	// Expand ~ in claude_config_dir
	if len(cfg.ClaudeConfigDir) > 0 && cfg.ClaudeConfigDir[0] == '~' {
		home, _ := os.UserHomeDir()
		cfg.ClaudeConfigDir = filepath.Join(home, cfg.ClaudeConfigDir[1:])
	}

	return cfg, nil
}

// ConfigPath returns the path to the config file.
func ConfigPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ConfigDir, ConfigFile)
}

// SkrDir returns the skr config directory.
func SkrDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ConfigDir)
}

// CachePath returns the path to the cached marketplace.json.
func CachePath() string {
	return filepath.Join(SkrDir(), CacheDir, "marketplace.json")
}

// LockPath returns the path to the lockfile.
func LockPath() string {
	return filepath.Join(SkrDir(), LockFile)
}

// Save writes the config to disk.
func (c *Config) Save() error {
	path := ConfigPath()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}
