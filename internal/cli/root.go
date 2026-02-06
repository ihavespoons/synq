package cli

import (
	"github.com/ihavespoons/synq/internal/config"
	"github.com/ihavespoons/synq/internal/logger"
	"github.com/spf13/cobra"
)

var (
	verbose   bool
	configDir string
)

func newRootCmd(version string) *cobra.Command {
	root := &cobra.Command{
		Use:     "synq",
		Short:   "Sync configuration files across machines via GitHub",
		Version: version,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			logger.Init(verbose)
		},
	}

	root.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "enable debug logging")
	root.PersistentFlags().StringVar(&configDir, "config-dir", config.DefaultConfigDir(), "synq config directory")

	root.AddCommand(
		newSetupCmd(),
		newAddCmd(),
		newRemoveCmd(),
		newListCmd(),
		newSyncCmd(),
		newDaemonCmd(),
	)

	return root
}

// Execute runs the root command.
func Execute(version string) error {
	return newRootCmd(version).Execute()
}
