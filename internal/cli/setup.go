package cli

import (
	"fmt"
	"os"

	"github.com/ihavespoons/synq/internal/config"
	"github.com/ihavespoons/synq/internal/daemon"
	"github.com/ihavespoons/synq/internal/fileops"
	"github.com/ihavespoons/synq/internal/gitops"
	"github.com/ihavespoons/synq/internal/logger"
	"github.com/spf13/cobra"
)

func newSetupCmd() *cobra.Command {
	var user string

	cmd := &cobra.Command{
		Use:   "setup",
		Short: "Initialize synq: create GitHub repo and clone locally",
		RunE: func(cmd *cobra.Command, args []string) error {
			log := logger.Get()

			// 1. Check gh CLI.
			log.Debug().Msg("checking gh CLI")
			if err := gitops.CheckGHInstalled(); err != nil {
				return err
			}
			fmt.Println("✓ gh CLI authenticated")

			// 2. Detect user.
			ghUser, err := gitops.DetectUser(user)
			if err != nil {
				return err
			}
			log.Debug().Str("user", ghUser).Msg("detected GitHub user")
			fmt.Printf("✓ GitHub user: %s\n", ghUser)

			// 3. Check/create repo.
			repoName := config.DefaultRepoName
			if gitops.RepoExists(ghUser, repoName) {
				fmt.Printf("✓ Repo %s/%s already exists\n", ghUser, repoName)
			} else {
				log.Debug().Msg("creating private repo")
				if err := gitops.CreatePrivateRepo(repoName); err != nil {
					return err
				}
				fmt.Printf("✓ Created private repo %s/%s\n", ghUser, repoName)
			}

			// 4. Clone repo.
			repoDir := config.RepoDir(configDir)
			if _, err := os.Stat(repoDir); err == nil {
				fmt.Printf("✓ Repo already cloned at %s\n", repoDir)
			} else {
				cloneURL, err := gitops.GetCloneURL(ghUser, repoName)
				if err != nil {
					return err
				}
				log.Debug().Str("url", cloneURL).Str("dir", repoDir).Msg("cloning repo")
				if err := gitops.Clone(cloneURL, repoDir); err != nil {
					return err
				}
				fmt.Printf("✓ Cloned to %s\n", repoDir)
			}

			// 5. Initialize repo config if needed.
			repoConfigPath := config.RepoConfigPath(configDir)
			if _, err := os.Stat(repoConfigPath); os.IsNotExist(err) {
				cfg := &config.Config{Files: []config.FileEntry{}}
				if err := config.SaveRepoConfig(configDir, cfg); err != nil {
					return fmt.Errorf("write repo config: %w", err)
				}
				if err := gitops.CommitAndPush(repoDir, "Initialize synq config"); err != nil {
					return fmt.Errorf("initial commit: %w", err)
				}
				fmt.Println("✓ Initialized synq.yaml in repo")
			}

			// 6. Write local state.
			cloneURL, _ := gitops.GetCloneURL(ghUser, repoName)
			state := &config.LocalState{
				GitHubUser: ghUser,
				RepoName:   repoName,
				RepoURL:    cloneURL,
				RepoPath:   fileops.TildePath(repoDir),
				Daemon: config.DaemonConfig{
					PollInterval: config.DefaultPollInterval,
				},
			}
			if err := config.SaveLocalState(configDir, state); err != nil {
				return fmt.Errorf("write local state: %w", err)
			}
			fmt.Println("✓ Saved local state")

			// 7. Install OS service.
			if err := daemon.InstallService(configDir); err != nil {
				log.Warn().Err(err).Msg("could not install service")
				fmt.Printf("⚠ Could not install OS service: %v\n", err)
			} else {
				fmt.Println("✓ Installed OS service")
			}

			fmt.Println("\nSetup complete! Use 'synq add <file>' to start managing files.")
			return nil
		},
	}

	cmd.Flags().StringVar(&user, "user", "", "GitHub username (auto-detected if omitted)")
	return cmd
}
