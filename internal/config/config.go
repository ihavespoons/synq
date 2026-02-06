package config

// Config is stored in the git repo as synq.yaml.
type Config struct {
	Files []FileEntry `yaml:"files"`
}

// FileEntry represents a single managed file.
type FileEntry struct {
	Name    string            `yaml:"name"`
	Source  string            `yaml:"source"`
	Targets map[string]string `yaml:"targets"`
}

// LocalState is stored at ~/.config/synq/synq.yaml.
type LocalState struct {
	GitHubUser string       `yaml:"github_user"`
	RepoName   string       `yaml:"repo_name"`
	RepoURL    string       `yaml:"repo_url"`
	RepoPath   string       `yaml:"repo_path"`
	Daemon     DaemonConfig `yaml:"daemon"`
}

// DaemonConfig holds daemon-specific settings.
type DaemonConfig struct {
	PollInterval string `yaml:"poll_interval"`
}

// SyncStatus represents the status of a managed file.
type SyncStatus string

const (
	StatusSynced   SyncStatus = "synced"
	StatusModified SyncStatus = "modified"
	StatusMissing  SyncStatus = "missing"
	StatusUnlinked SyncStatus = "unlinked"
	StatusNoTarget SyncStatus = "no target"
)
