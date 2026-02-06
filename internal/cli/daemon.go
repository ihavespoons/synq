package cli

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/ihavespoons/synq/internal/daemon"
	"github.com/ihavespoons/synq/internal/logger"
	"github.com/spf13/cobra"
)

func newDaemonCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "daemon",
		Short: "Manage the synq background daemon",
	}

	cmd.AddCommand(
		newDaemonStartCmd(),
		newDaemonStopCmd(),
		newDaemonStatusCmd(),
	)

	return cmd
}

func newDaemonStartCmd() *cobra.Command {
	var foreground bool

	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start the synq daemon",
		RunE: func(cmd *cobra.Command, args []string) error {
			log := logger.Get()

			if foreground {
				// Run the main loop directly.
				return daemon.Run(configDir)
			}

			// Check if already running.
			if running, pid := daemon.IsRunning(configDir); running {
				fmt.Printf("Daemon already running (PID %d)\n", pid)
				return nil
			}

			// Start self in background with --foreground.
			exe, err := os.Executable()
			if err != nil {
				return fmt.Errorf("find executable: %w", err)
			}

			bgArgs := []string{"daemon", "start", "--foreground", "--config-dir", configDir}
			if verbose {
				bgArgs = append(bgArgs, "--verbose")
			}

			bgCmd := exec.Command(exe, bgArgs...)
			bgCmd.Stdout = nil
			bgCmd.Stderr = nil
			bgCmd.Stdin = nil

			// Detach from parent.
			bgCmd.SysProcAttr = nil

			if err := bgCmd.Start(); err != nil {
				return fmt.Errorf("start daemon: %w", err)
			}

			// Release the process so it continues after we exit.
			if err := bgCmd.Process.Release(); err != nil {
				log.Warn().Err(err).Msg("release process")
			}

			fmt.Printf("Daemon started (PID %d)\n", bgCmd.Process.Pid)
			return nil
		},
	}

	cmd.Flags().BoolVar(&foreground, "foreground", false, "run in foreground (used internally)")
	return cmd
}

func newDaemonStopCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "stop",
		Short: "Stop the synq daemon",
		RunE: func(cmd *cobra.Command, args []string) error {
			running, pid := daemon.IsRunning(configDir)
			if !running {
				fmt.Println("Daemon is not running")
				return nil
			}

			if err := daemon.StopProcess(pid); err != nil {
				return fmt.Errorf("stop daemon: %w", err)
			}

			daemon.RemovePID(configDir)
			fmt.Printf("Daemon stopped (PID %d)\n", pid)
			return nil
		},
	}
}

func newDaemonStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Check if the synq daemon is running",
		RunE: func(cmd *cobra.Command, args []string) error {
			running, pid := daemon.IsRunning(configDir)
			if running {
				fmt.Printf("Daemon is running (PID %d)\n", pid)
			} else {
				fmt.Println("Daemon is not running")
			}
			return nil
		},
	}
}
