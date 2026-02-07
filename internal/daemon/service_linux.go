//go:build linux

package daemon

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

const systemdUnit = `[Unit]
Description=synq configuration sync daemon

[Service]
ExecStart=%s daemon start --foreground
Restart=on-failure
RestartSec=10

[Install]
WantedBy=default.target
`

func unitPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "systemd", "user", "synq.service")
}

// InstallService installs the systemd user unit for Linux.
func InstallService(configDir string) error {
	binPath, err := exec.LookPath("synq")
	if err != nil {
		binPath, err = os.Executable()
		if err != nil {
			return fmt.Errorf("find synq binary: %w", err)
		}
	}

	path := unitPath()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	content := fmt.Sprintf(systemdUnit, binPath)
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		return err
	}

	// Reload systemd, enable and start the service.
	_ = exec.Command("systemctl", "--user", "daemon-reload").Run()
	_ = exec.Command("systemctl", "--user", "enable", "synq").Run()
	if err := exec.Command("systemctl", "--user", "start", "synq").Run(); err != nil {
		return fmt.Errorf("start service: %w (you may need to run: loginctl enable-linger $USER)", err)
	}
	return nil
}

// UninstallService removes the systemd user unit.
func UninstallService() error {
	_ = exec.Command("systemctl", "--user", "stop", "synq").Run()
	_ = exec.Command("systemctl", "--user", "disable", "synq").Run()
	return os.Remove(unitPath())
}
