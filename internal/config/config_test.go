package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSaveAndLoadLocalState(t *testing.T) {
	tmp := t.TempDir()

	state := &LocalState{
		GitHubUser: "testuser",
		RepoName:   "synq-config",
		RepoURL:    "git@github.com:testuser/synq-config.git",
		RepoPath:   "~/.config/synq/repo",
		Daemon: DaemonConfig{
			PollInterval: "5m",
		},
	}

	if err := SaveLocalState(tmp, state); err != nil {
		t.Fatal(err)
	}

	loaded, err := LoadLocalState(tmp)
	if err != nil {
		t.Fatal(err)
	}

	if loaded.GitHubUser != state.GitHubUser {
		t.Errorf("GitHubUser = %q, want %q", loaded.GitHubUser, state.GitHubUser)
	}
	if loaded.RepoName != state.RepoName {
		t.Errorf("RepoName = %q, want %q", loaded.RepoName, state.RepoName)
	}
	if loaded.RepoURL != state.RepoURL {
		t.Errorf("RepoURL = %q, want %q", loaded.RepoURL, state.RepoURL)
	}
	if loaded.Daemon.PollInterval != state.Daemon.PollInterval {
		t.Errorf("PollInterval = %q, want %q", loaded.Daemon.PollInterval, state.Daemon.PollInterval)
	}
}

func TestSaveAndLoadRepoConfig(t *testing.T) {
	tmp := t.TempDir()
	// Create the "repo" subdirectory since RepoConfigPath expects it.
	repoDir := filepath.Join(tmp, "repo")
	if err := os.MkdirAll(repoDir, 0o755); err != nil {
		t.Fatal(err)
	}

	cfg := &Config{
		Files: []FileEntry{
			{
				Name:   "starship.toml",
				Source: "starship.toml",
				Targets: map[string]string{
					"darwin": "~/.config/starship.toml",
					"linux":  "~/.config/starship.toml",
				},
			},
		},
	}

	if err := SaveRepoConfig(tmp, cfg); err != nil {
		t.Fatal(err)
	}

	loaded, err := LoadRepoConfig(tmp)
	if err != nil {
		t.Fatal(err)
	}

	if len(loaded.Files) != 1 {
		t.Fatalf("expected 1 file, got %d", len(loaded.Files))
	}
	if loaded.Files[0].Name != "starship.toml" {
		t.Errorf("Name = %q, want %q", loaded.Files[0].Name, "starship.toml")
	}
	if loaded.Files[0].Targets["darwin"] != "~/.config/starship.toml" {
		t.Errorf("target = %q", loaded.Files[0].Targets["darwin"])
	}
}

func TestLoadRepoConfig_NotExist(t *testing.T) {
	tmp := t.TempDir()
	cfg, err := LoadRepoConfig(tmp)
	if err != nil {
		t.Fatal(err)
	}
	if len(cfg.Files) != 0 {
		t.Errorf("expected empty config, got %d files", len(cfg.Files))
	}
}

func TestDefaultConfigDir(t *testing.T) {
	dir := DefaultConfigDir()
	if dir == "" {
		t.Error("expected non-empty config dir")
	}
}
