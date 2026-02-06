package cli

import (
	"fmt"
	"path/filepath"

	"github.com/ihavespoons/synq/internal/config"
	"github.com/ihavespoons/synq/internal/fileops"
	"github.com/ihavespoons/synq/internal/gitops"
	"github.com/ihavespoons/synq/internal/logger"
	"github.com/spf13/cobra"
)

func newSyncCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "sync",
		Short: "Sync configuration files with remote repo",
		RunE: func(cmd *cobra.Command, args []string) error {
			log := logger.Get()
			repoDir := config.RepoDir(configDir)

			// 1. Commit + push local changes.
			if gitops.HasChanges(repoDir) {
				log.Debug().Msg("committing local changes")
				if err := gitops.CommitAndPush(repoDir, "Sync local changes"); err != nil {
					return fmt.Errorf("push local changes: %w", err)
				}
				fmt.Println("✓ Pushed local changes")
			}

			// 2. Pull remote changes.
			log.Debug().Msg("pulling remote changes")
			changed, err := gitops.Pull(repoDir)
			if err != nil {
				return fmt.Errorf("pull: %w", err)
			}
			if changed {
				fmt.Println("✓ Pulled remote changes")
			} else {
				fmt.Println("✓ Already up to date")
			}

			// 3. Re-apply symlinks.
			cfg, err := config.LoadRepoConfig(configDir)
			if err != nil {
				return fmt.Errorf("load repo config: %w", err)
			}

			for _, f := range cfg.Files {
				target, ok := fileops.ResolveTarget(f.Targets)
				if !ok {
					log.Debug().Str("name", f.Name).Msg("no target for this OS, skipping")
					continue
				}

				repoFile := filepath.Join(repoDir, f.Source)
				if fileops.IsSymlinkTo(target, repoFile) {
					continue
				}

				log.Debug().Str("name", f.Name).Str("target", target).Msg("creating symlink")
				if err := fileops.CreateSymlink(repoFile, target); err != nil {
					log.Error().Err(err).Str("name", f.Name).Msg("failed to create symlink")
					continue
				}
				fmt.Printf("✓ Linked %s -> %s\n", f.Name, fileops.TildePath(target))
			}

			return nil
		},
	}
}
