//go:build windows

package daemon

import (
	"os"
)

// processExists checks if a process with the given PID exists.
func processExists(pid int) bool {
	proc, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	// On Windows, FindProcess always succeeds; try to signal.
	err = proc.Signal(os.Signal(nil))
	return err == nil
}

// StopProcess kills the process on Windows.
func StopProcess(pid int) error {
	proc, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	return proc.Kill()
}
