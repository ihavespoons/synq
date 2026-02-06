package daemon

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/ihavespoons/synq/internal/config"
	"github.com/ihavespoons/synq/internal/fileops"
	"github.com/ihavespoons/synq/internal/gitops"
	"github.com/ihavespoons/synq/internal/logger"
)

// Run starts the daemon main loop. It blocks until a signal is received.
func Run(configDir string) error {
	log := logger.Get()
	log.Info().Msg("synq daemon starting")

	// Write PID file.
	if err := WritePID(configDir); err != nil {
		return fmt.Errorf("write PID: %w", err)
	}
	defer RemovePID(configDir)

	state, err := config.LoadLocalState(configDir)
	if err != nil {
		return fmt.Errorf("load local state: %w", err)
	}

	repoDir := config.RepoDir(configDir)
	pollInterval, err := time.ParseDuration(state.Daemon.PollInterval)
	if err != nil {
		pollInterval = 5 * time.Minute
	}

	// Set up file watcher.
	onChange := func() {
		log.Info().Msg("file change detected, syncing")
		if err := gitops.CommitAndPush(repoDir, "Auto-sync: file changed"); err != nil {
			log.Error().Err(err).Msg("auto-sync failed")
		}
	}

	watcher, err := NewWatcher(onChange, log)
	if err != nil {
		return fmt.Errorf("create watcher: %w", err)
	}
	defer watcher.Close()

	// Initial watch setup.
	refreshWatcher(configDir, watcher)
	watcher.Start()

	// Poll ticker.
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	// Signal handling.
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	log.Info().Str("poll_interval", pollInterval.String()).Msg("daemon running")

	for {
		select {
		case <-ticker.C:
			log.Debug().Msg("poll tick: pulling changes")
			changed, err := gitops.Pull(repoDir)
			if err != nil {
				log.Error().Err(err).Msg("pull failed")
				continue
			}
			if changed {
				log.Info().Msg("remote changes found, applying symlinks")
				applySymlinks(configDir)
				refreshWatcher(configDir, watcher)
			}

		case sig := <-sigCh:
			log.Info().Str("signal", sig.String()).Msg("shutting down")
			return nil
		}
	}
}

func refreshWatcher(configDir string, watcher *Watcher) {
	log := logger.Get()
	cfg, err := config.LoadRepoConfig(configDir)
	if err != nil {
		log.Error().Err(err).Msg("load repo config for watcher")
		return
	}

	var paths []string
	repoDir := config.RepoDir(configDir)
	for _, f := range cfg.Files {
		target, ok := fileops.ResolveTarget(f.Targets)
		if ok {
			paths = append(paths, target)
		}
		paths = append(paths, filepath.Join(repoDir, f.Source))
	}
	watcher.WatchPaths(paths)
}

func applySymlinks(configDir string) {
	log := logger.Get()
	cfg, err := config.LoadRepoConfig(configDir)
	if err != nil {
		log.Error().Err(err).Msg("load repo config for symlinks")
		return
	}

	repoDir := config.RepoDir(configDir)
	for _, f := range cfg.Files {
		target, ok := fileops.ResolveTarget(f.Targets)
		if !ok {
			continue
		}
		repoFile := filepath.Join(repoDir, f.Source)
		if fileops.IsSymlinkTo(target, repoFile) {
			continue
		}

		// Handle conflict: backup local file.
		if info, err := os.Lstat(target); err == nil && info.Mode()&os.ModeSymlink == 0 {
			backup := target + fmt.Sprintf(".conflict-%s", time.Now().Format("20060102-150405"))
			if err := os.Rename(target, backup); err != nil {
				log.Error().Err(err).Str("file", target).Msg("backup failed")
				continue
			}
			log.Info().Str("backup", backup).Msg("backed up conflicting file")
		}

		if err := fileops.CreateSymlink(repoFile, target); err != nil {
			log.Error().Err(err).Str("name", f.Name).Msg("symlink failed")
		}
	}
}
