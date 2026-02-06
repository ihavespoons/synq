package fileops

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCopyFile(t *testing.T) {
	tmp := t.TempDir()
	src := filepath.Join(tmp, "src.txt")
	dst := filepath.Join(tmp, "subdir", "dst.txt")

	if err := os.WriteFile(src, []byte("hello"), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := CopyFile(src, dst); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(dst)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "hello" {
		t.Errorf("copied content = %q, want %q", string(data), "hello")
	}
}

func TestCreateSymlink(t *testing.T) {
	tmp := t.TempDir()
	target := filepath.Join(tmp, "target.txt")
	link := filepath.Join(tmp, "link.txt")

	if err := os.WriteFile(target, []byte("content"), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := CreateSymlink(target, link); err != nil {
		t.Fatal(err)
	}

	dest, err := os.Readlink(link)
	if err != nil {
		t.Fatal(err)
	}
	if dest != target {
		t.Errorf("symlink points to %q, want %q", dest, target)
	}

	// Verify we can read through the symlink.
	data, err := os.ReadFile(link)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "content" {
		t.Errorf("read through symlink = %q, want %q", string(data), "content")
	}
}

func TestCreateSymlink_ReplacesExisting(t *testing.T) {
	tmp := t.TempDir()
	target := filepath.Join(tmp, "target.txt")
	link := filepath.Join(tmp, "link.txt")

	if err := os.WriteFile(target, []byte("content"), 0o644); err != nil {
		t.Fatal(err)
	}
	// Create existing file at link location.
	if err := os.WriteFile(link, []byte("old"), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := CreateSymlink(target, link); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(link)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "content" {
		t.Errorf("read = %q, want %q", string(data), "content")
	}
}

func TestIsSymlinkTo(t *testing.T) {
	tmp := t.TempDir()
	target := filepath.Join(tmp, "target.txt")
	link := filepath.Join(tmp, "link.txt")

	if err := os.WriteFile(target, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(target, link); err != nil {
		t.Fatal(err)
	}

	if !IsSymlinkTo(link, target) {
		t.Error("expected IsSymlinkTo to return true")
	}

	// Wrong target.
	other := filepath.Join(tmp, "other.txt")
	if IsSymlinkTo(link, other) {
		t.Error("expected IsSymlinkTo to return false for wrong target")
	}

	// Not a symlink.
	if IsSymlinkTo(target, target) {
		t.Error("expected IsSymlinkTo to return false for regular file")
	}
}

func TestRemoveSymlink(t *testing.T) {
	tmp := t.TempDir()
	repoFile := filepath.Join(tmp, "repo", "test.txt")
	linkPath := filepath.Join(tmp, "link.txt")

	// Create repo file.
	if err := os.MkdirAll(filepath.Dir(repoFile), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(repoFile, []byte("repo content"), 0o644); err != nil {
		t.Fatal(err)
	}

	// Create symlink.
	if err := os.Symlink(repoFile, linkPath); err != nil {
		t.Fatal(err)
	}

	// Remove symlink and restore.
	if err := RemoveSymlink(linkPath, repoFile); err != nil {
		t.Fatal(err)
	}

	// Should be a regular file now.
	info, err := os.Lstat(linkPath)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode()&os.ModeSymlink != 0 {
		t.Error("expected regular file, got symlink")
	}

	data, err := os.ReadFile(linkPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "repo content" {
		t.Errorf("restored content = %q, want %q", string(data), "repo content")
	}
}
