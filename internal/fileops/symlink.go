package fileops

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// CopyFile copies src to dst, creating parent directories as needed.
func CopyFile(src, dst string) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return fmt.Errorf("create parent dir: %w", err)
	}

	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func() { _ = in.Close() }()

	info, err := in.Stat()
	if err != nil {
		return err
	}

	out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, info.Mode())
	if err != nil {
		return err
	}
	defer func() { _ = out.Close() }()

	_, err = io.Copy(out, in)
	return err
}

// CreateSymlink creates a symlink at linkPath pointing to target.
// It removes any existing file at linkPath first.
func CreateSymlink(target, linkPath string) error {
	if err := os.MkdirAll(filepath.Dir(linkPath), 0o755); err != nil {
		return fmt.Errorf("create parent dir: %w", err)
	}

	// Remove existing file/symlink.
	if _, err := os.Lstat(linkPath); err == nil {
		if err := os.Remove(linkPath); err != nil {
			return fmt.Errorf("remove existing file: %w", err)
		}
	}

	return os.Symlink(target, linkPath)
}

// RemoveSymlink removes a symlink and copies the target file back to the link location.
func RemoveSymlink(linkPath, repoFilePath string) error {
	info, err := os.Lstat(linkPath)
	if err != nil {
		// Link doesn't exist; just copy from repo if the repo file exists.
		if os.IsNotExist(err) {
			return CopyFile(repoFilePath, linkPath)
		}
		return err
	}

	if info.Mode()&os.ModeSymlink == 0 {
		// Not a symlink; leave it alone.
		return nil
	}

	if err := os.Remove(linkPath); err != nil {
		return fmt.Errorf("remove symlink: %w", err)
	}

	return CopyFile(repoFilePath, linkPath)
}

// IsSymlinkTo returns true if path is a symlink pointing to target.
func IsSymlinkTo(path, target string) bool {
	dest, err := os.Readlink(path)
	if err != nil {
		return false
	}
	// Resolve both to absolute for comparison.
	absTarget, err1 := filepath.Abs(target)
	absDest, err2 := filepath.Abs(dest)
	if err1 != nil || err2 != nil {
		return dest == target
	}
	return absTarget == absDest
}
