package daemon

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// PIDFile returns the path to the daemon PID file.
func PIDFile(configDir string) string {
	return filepath.Join(configDir, "daemon.pid")
}

// WritePID writes the current process ID to the PID file.
func WritePID(configDir string) error {
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		return err
	}
	return os.WriteFile(PIDFile(configDir), []byte(strconv.Itoa(os.Getpid())), 0o644)
}

// RemovePID removes the PID file.
func RemovePID(configDir string) {
	os.Remove(PIDFile(configDir))
}

// ReadPID reads the PID from the PID file. Returns 0 if not found.
func ReadPID(configDir string) (int, error) {
	data, err := os.ReadFile(PIDFile(configDir))
	if err != nil {
		if os.IsNotExist(err) {
			return 0, nil
		}
		return 0, err
	}
	pid, err := strconv.Atoi(strings.TrimSpace(string(data)))
	if err != nil {
		return 0, fmt.Errorf("invalid PID file: %w", err)
	}
	return pid, nil
}

// IsRunning checks if the daemon is running by reading the PID file
// and checking if the process exists.
func IsRunning(configDir string) (bool, int) {
	pid, err := ReadPID(configDir)
	if err != nil || pid == 0 {
		return false, 0
	}
	if !processExists(pid) {
		// Stale PID file.
		RemovePID(configDir)
		return false, 0
	}
	return true, pid
}
