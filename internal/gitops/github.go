package gitops

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

// CheckGHInstalled verifies the gh CLI is available and authenticated.
func CheckGHInstalled() error {
	if _, err := exec.LookPath("gh"); err != nil {
		return fmt.Errorf("gh CLI not found; install from https://cli.github.com")
	}
	out, err := exec.Command("gh", "auth", "status").CombinedOutput()
	if err != nil {
		return fmt.Errorf("gh not authenticated: %s", strings.TrimSpace(string(out)))
	}
	return nil
}

// DetectUser returns the GitHub username.
// Priority: explicit flag > gh api > git config.
func DetectUser(explicit string) (string, error) {
	if explicit != "" {
		return explicit, nil
	}
	out, err := exec.Command("gh", "api", "user", "--jq", ".login").Output()
	if err == nil {
		user := strings.TrimSpace(string(out))
		if user != "" {
			return user, nil
		}
	}
	out, err = exec.Command("git", "config", "user.name").Output()
	if err == nil {
		user := strings.TrimSpace(string(out))
		if user != "" {
			return user, nil
		}
	}
	return "", fmt.Errorf("could not detect GitHub username; use --user flag")
}

// RepoExists checks if the named repo exists for the user.
func RepoExists(user, repo string) bool {
	err := exec.Command("gh", "repo", "view", user+"/"+repo, "--json", "name").Run()
	return err == nil
}

// CreatePrivateRepo creates a private repo on GitHub.
func CreatePrivateRepo(repo string) error {
	out, err := exec.Command("gh", "repo", "create", repo,
		"--private",
		"--description", "synq configuration files",
	).CombinedOutput()
	if err != nil {
		return fmt.Errorf("create repo: %s", strings.TrimSpace(string(out)))
	}
	return nil
}

// GetCloneURL returns the SSH clone URL for a repo.
func GetCloneURL(user, repo string) (string, error) {
	out, err := exec.Command("gh", "repo", "view", user+"/"+repo, "--json", "sshUrl").Output()
	if err != nil {
		return "", fmt.Errorf("get clone URL: %w", err)
	}
	var result struct {
		SSHURL string `json:"sshUrl"`
	}
	if err := json.Unmarshal(out, &result); err != nil {
		return "", err
	}
	return result.SSHURL, nil
}
