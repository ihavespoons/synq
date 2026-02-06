package config

import (
	"os"
	"path/filepath"

	"github.com/ihavespoons/synq/internal/fileops"
	"gopkg.in/yaml.v3"
)

const (
	DefaultRepoName     = "synq-config"
	DefaultPollInterval = "5m"
	RepoConfigFile      = "synq.yaml"
)

// DefaultConfigDir returns ~/.config/synq.
func DefaultConfigDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "synq")
}

// LocalStatePath returns the path to the local state file.
func LocalStatePath(configDir string) string {
	return filepath.Join(configDir, "synq.yaml")
}

// RepoDir returns the path to the cloned repo.
func RepoDir(configDir string) string {
	return filepath.Join(configDir, "repo")
}

// RepoConfigPath returns the path to synq.yaml inside the repo.
func RepoConfigPath(configDir string) string {
	return filepath.Join(RepoDir(configDir), RepoConfigFile)
}

// LoadLocalState reads the local state file.
func LoadLocalState(configDir string) (*LocalState, error) {
	data, err := os.ReadFile(LocalStatePath(configDir))
	if err != nil {
		return nil, err
	}
	var state LocalState
	if err := yaml.Unmarshal(data, &state); err != nil {
		return nil, err
	}
	return &state, nil
}

// SaveLocalState writes the local state file.
func SaveLocalState(configDir string, state *LocalState) error {
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		return err
	}
	data, err := yaml.Marshal(state)
	if err != nil {
		return err
	}
	return os.WriteFile(LocalStatePath(configDir), data, 0o644)
}

// LoadRepoConfig reads synq.yaml from the cloned repo.
func LoadRepoConfig(configDir string) (*Config, error) {
	path := RepoConfigPath(configDir)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{}, nil
		}
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// SaveRepoConfig writes synq.yaml to the cloned repo.
func SaveRepoConfig(configDir string, cfg *Config) error {
	path := RepoConfigPath(configDir)
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

// ExpandRepoPath expands the repo_path field from local state.
func ExpandRepoPath(state *LocalState) string {
	return fileops.ExpandPath(state.RepoPath)
}
