package fileops

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestExpandPath_Tilde(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Skip("no home dir")
	}
	got := ExpandPath("~/foo/bar")
	want := filepath.Join(home, "foo", "bar")
	if got != want {
		t.Errorf("ExpandPath(~/foo/bar) = %q, want %q", got, want)
	}
}

func TestExpandPath_Absolute(t *testing.T) {
	got := ExpandPath("/tmp/test")
	if got != filepath.Clean("/tmp/test") {
		t.Errorf("ExpandPath(/tmp/test) = %q", got)
	}
}

func TestCurrentOSKey(t *testing.T) {
	key := CurrentOSKey()
	if key != runtime.GOOS {
		t.Errorf("CurrentOSKey() = %q, want %q", key, runtime.GOOS)
	}
}

func TestResolveTarget(t *testing.T) {
	targets := map[string]string{
		"darwin": "~/.config/test",
		"linux":  "~/.config/test",
	}
	path, ok := ResolveTarget(targets)
	if runtime.GOOS == "darwin" || runtime.GOOS == "linux" {
		if !ok {
			t.Error("expected target to be found")
		}
		if path == "" {
			t.Error("expected non-empty path")
		}
	}
}

func TestResolveTarget_Missing(t *testing.T) {
	targets := map[string]string{
		"nonexistent_os": "/some/path",
	}
	_, ok := ResolveTarget(targets)
	if ok {
		t.Error("expected no target for nonexistent OS")
	}
}

func TestTildePath(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Skip("no home dir")
	}
	got := TildePath(filepath.Join(home, "foo", "bar"))
	if got != "~/foo/bar" {
		t.Errorf("TildePath() = %q, want %q", got, "~/foo/bar")
	}
}

func TestTildePath_NonHome(t *testing.T) {
	got := TildePath("/tmp/something")
	if got != "/tmp/something" {
		t.Errorf("TildePath(/tmp/something) = %q", got)
	}
}
