package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"text/tabwriter"

	"github.com/ihavespoons/synq/internal/config"
	"github.com/ihavespoons/synq/internal/fileops"
	"github.com/spf13/cobra"
)

func newListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all managed files and their status",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.LoadRepoConfig(configDir)
			if err != nil {
				return fmt.Errorf("load repo config: %w", err)
			}

			if len(cfg.Files) == 0 {
				fmt.Println("No files managed by synq. Use 'synq add <file>' to get started.")
				return nil
			}

			repoDir := config.RepoDir(configDir)

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "NAME\tTARGET\tSTATUS")

			for _, f := range cfg.Files {
				target, hasTarget := fileops.ResolveTarget(f.Targets)
				status := getStatus(f, repoDir, target, hasTarget)

				targetDisplay := "(no target for this OS)"
				if hasTarget {
					targetDisplay = fileops.TildePath(target)
				}

				fmt.Fprintf(w, "%s\t%s\t%s\n", f.Name, targetDisplay, status)
			}

			return w.Flush()
		},
	}
}

func getStatus(f config.FileEntry, repoDir, target string, hasTarget bool) config.SyncStatus {
	if !hasTarget {
		return config.StatusNoTarget
	}

	repoFile := filepath.Join(repoDir, f.Source)

	// Check if repo file exists.
	if _, err := os.Stat(repoFile); os.IsNotExist(err) {
		return config.StatusMissing
	}

	// Check if target exists.
	info, err := os.Lstat(target)
	if os.IsNotExist(err) {
		return config.StatusMissing
	}
	if err != nil {
		return config.StatusMissing
	}

	// Check if it's a symlink pointing to the repo file.
	if info.Mode()&os.ModeSymlink == 0 {
		return config.StatusUnlinked
	}

	if !fileops.IsSymlinkTo(target, repoFile) {
		return config.StatusUnlinked
	}

	// Check if repo has uncommitted changes for this file.
	return config.StatusSynced
}
