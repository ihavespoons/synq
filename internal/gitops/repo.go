package gitops

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func git(dir string, args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(out)), err
}

// Clone clones a repo to the given directory.
func Clone(url, dir string) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	out, err := exec.Command("git", "clone", url, dir).CombinedOutput()
	if err != nil {
		return fmt.Errorf("git clone: %s", strings.TrimSpace(string(out)))
	}
	return nil
}

// Add stages files in the repo.
func Add(repoDir string, files ...string) error {
	args := append([]string{"add"}, files...)
	if out, err := git(repoDir, args...); err != nil {
		return fmt.Errorf("git add: %s", out)
	}
	return nil
}

// AddAll stages all changes.
func AddAll(repoDir string) error {
	if out, err := git(repoDir, "add", "-A"); err != nil {
		return fmt.Errorf("git add -A: %s", out)
	}
	return nil
}

// Commit creates a commit with the given message.
// Returns false if there was nothing to commit.
func Commit(repoDir, message string) (bool, error) {
	out, err := git(repoDir, "commit", "-m", message)
	if err != nil {
		if strings.Contains(out, "nothing to commit") {
			return false, nil
		}
		return false, fmt.Errorf("git commit: %s", out)
	}
	return true, nil
}

// Push pushes to origin.
func Push(repoDir string) error {
	out, err := git(repoDir, "push")
	if err != nil {
		return fmt.Errorf("git push: %s", out)
	}
	return nil
}

// Pull pulls from origin. Returns true if new changes were fetched.
func Pull(repoDir string) (bool, error) {
	out, err := git(repoDir, "pull", "--rebase")
	if err != nil {
		return false, fmt.Errorf("git pull: %s", out)
	}
	return !strings.Contains(out, "Already up to date"), nil
}

// HasChanges returns true if there are uncommitted changes.
func HasChanges(repoDir string) bool {
	out, _ := git(repoDir, "status", "--porcelain")
	return out != ""
}

// CommitAndPush stages all, commits, and pushes.
func CommitAndPush(repoDir, message string) error {
	if err := AddAll(repoDir); err != nil {
		return err
	}
	committed, err := Commit(repoDir, message)
	if err != nil {
		return err
	}
	if !committed {
		return nil
	}
	return Push(repoDir)
}

// InitRepo initializes a new git repo if the directory is not already one.
func InitRepo(dir string) error {
	if out, err := git(dir, "rev-parse", "--is-inside-work-tree"); err != nil || out != "true" {
		if out, err := git(dir, "init"); err != nil {
			return fmt.Errorf("git init: %s", out)
		}
	}
	return nil
}
