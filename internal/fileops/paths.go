package fileops

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// ExpandPath expands ~ and environment variables in a path.
func ExpandPath(path string) string {
	if strings.HasPrefix(path, "~/") || path == "~" {
		home, err := os.UserHomeDir()
		if err == nil {
			path = filepath.Join(home, path[1:])
		}
	}
	// Expand %APPDATA% and similar on Windows.
	if runtime.GOOS == "windows" {
		path = os.ExpandEnv(path)
	}
	return filepath.Clean(path)
}

// CurrentOSKey returns the OS key used in target maps.
func CurrentOSKey() string {
	return runtime.GOOS
}

// ResolveTarget returns the expanded target path for the current OS.
func ResolveTarget(targets map[string]string) (string, bool) {
	t, ok := targets[CurrentOSKey()]
	if !ok {
		return "", false
	}
	return ExpandPath(t), true
}

// TildePath collapses a home-relative path back to ~/...
func TildePath(abs string) string {
	home, err := os.UserHomeDir()
	if err != nil {
		return abs
	}
	if rel, err := filepath.Rel(home, abs); err == nil && !strings.HasPrefix(rel, "..") {
		return "~/" + rel
	}
	return abs
}
