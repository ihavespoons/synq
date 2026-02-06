package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/ihavespoons/synq/internal/config"
	"github.com/ihavespoons/synq/internal/fileops"
	"github.com/ihavespoons/synq/internal/gitops"
	"github.com/ihavespoons/synq/internal/logger"
	"github.com/spf13/cobra"
)

func newAddCmd() *cobra.Command {
	var name string

	cmd := &cobra.Command{
		Use:   "add <file>",
		Short: "Add a file to synq management",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			log := logger.Get()

			// 1. Resolve file path.
			filePath := fileops.ExpandPath(args[0])
			absPath, err := filepath.Abs(filePath)
			if err != nil {
				return fmt.Errorf("resolve path: %w", err)
			}
			if _, err := os.Stat(absPath); err != nil {
				return fmt.Errorf("file not found: %s", absPath)
			}

			// Determine name.
			if name == "" {
				name = filepath.Base(absPath)
			}
			log.Debug().Str("file", absPath).Str("name", name).Msg("adding file")

			repoDir := config.RepoDir(configDir)
			repoFilePath := filepath.Join(repoDir, name)

			// 2. Copy file into repo.
			if err := fileops.CopyFile(absPath, repoFilePath); err != nil {
				return fmt.Errorf("copy to repo: %w", err)
			}
			fmt.Printf("✓ Copied %s to repo as %s\n", filepath.Base(absPath), name)

			// 3. Replace original with symlink.
			if err := fileops.CreateSymlink(repoFilePath, absPath); err != nil {
				return fmt.Errorf("create symlink: %w", err)
			}
			fmt.Printf("✓ Created symlink %s -> %s\n", fileops.TildePath(absPath), name)

			// 4. Update repo config.
			cfg, err := config.LoadRepoConfig(configDir)
			if err != nil {
				return fmt.Errorf("load repo config: %w", err)
			}

			// Check if entry already exists.
			found := false
			for i, f := range cfg.Files {
				if f.Name == name {
					cfg.Files[i].Targets[runtime.GOOS] = fileops.TildePath(absPath)
					found = true
					break
				}
			}
			if !found {
				cfg.Files = append(cfg.Files, config.FileEntry{
					Name:   name,
					Source: name,
					Targets: map[string]string{
						runtime.GOOS: fileops.TildePath(absPath),
					},
				})
			}

			if err := config.SaveRepoConfig(configDir, cfg); err != nil {
				return fmt.Errorf("save repo config: %w", err)
			}

			// 5. Commit + push.
			if err := gitops.CommitAndPush(repoDir, fmt.Sprintf("Add %s", name)); err != nil {
				return fmt.Errorf("commit and push: %w", err)
			}
			fmt.Printf("✓ Committed and pushed\n")

			return nil
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "name for the file in the repo (defaults to filename)")
	return cmd
}
