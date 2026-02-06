package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ihavespoons/synq/internal/config"
	"github.com/ihavespoons/synq/internal/fileops"
	"github.com/ihavespoons/synq/internal/gitops"
	"github.com/ihavespoons/synq/internal/logger"
	"github.com/spf13/cobra"
)

func newRemoveCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "remove <name>",
		Short: "Remove a file from synq management",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			log := logger.Get()
			name := args[0]

			// 1. Load config, find entry.
			cfg, err := config.LoadRepoConfig(configDir)
			if err != nil {
				return fmt.Errorf("load repo config: %w", err)
			}

			idx := -1
			for i, f := range cfg.Files {
				if f.Name == name {
					idx = i
					break
				}
			}
			if idx == -1 {
				return fmt.Errorf("file %q not found in synq config", name)
			}

			entry := cfg.Files[idx]
			repoDir := config.RepoDir(configDir)
			repoFilePath := filepath.Join(repoDir, entry.Source)

			// 2. Restore original file from symlink.
			targetPath, hasTarget := fileops.ResolveTarget(entry.Targets)
			if hasTarget {
				log.Debug().Str("target", targetPath).Msg("restoring original file")
				if err := fileops.RemoveSymlink(targetPath, repoFilePath); err != nil {
					return fmt.Errorf("restore file: %w", err)
				}
				fmt.Printf("✓ Restored %s\n", fileops.TildePath(targetPath))
			}

			// 3. Delete file from repo.
			if err := os.Remove(repoFilePath); err != nil && !os.IsNotExist(err) {
				return fmt.Errorf("remove repo file: %w", err)
			}

			// 4. Remove entry from config.
			cfg.Files = append(cfg.Files[:idx], cfg.Files[idx+1:]...)
			if err := config.SaveRepoConfig(configDir, cfg); err != nil {
				return fmt.Errorf("save repo config: %w", err)
			}

			// 5. Commit + push.
			if err := gitops.CommitAndPush(repoDir, fmt.Sprintf("Remove %s", name)); err != nil {
				return fmt.Errorf("commit and push: %w", err)
			}
			fmt.Printf("✓ Removed %s from synq\n", name)

			return nil
		},
	}
}
