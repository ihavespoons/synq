//go:build darwin

package daemon

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"
)

const plistTemplate = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>com.synq.daemon</string>
    <key>ProgramArguments</key>
    <array>
        <string>{{.BinaryPath}}</string>
        <string>daemon</string>
        <string>start</string>
        <string>--foreground</string>
    </array>
    <key>RunAtLoad</key>
    <true/>
    <key>KeepAlive</key>
    <true/>
    <key>StandardOutPath</key>
    <string>{{.LogDir}}/synq.log</string>
    <key>StandardErrorPath</key>
    <string>{{.LogDir}}/synq.err</string>
</dict>
</plist>
`

func plistPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, "Library", "LaunchAgents", "com.synq.daemon.plist")
}

// InstallService installs the launchd plist for macOS.
func InstallService(configDir string) error {
	binPath, err := exec.LookPath("synq")
	if err != nil {
		// Fall back to the current executable.
		binPath, err = os.Executable()
		if err != nil {
			return fmt.Errorf("find synq binary: %w", err)
		}
	}

	logDir := filepath.Join(configDir, "logs")
	if err := os.MkdirAll(logDir, 0o755); err != nil {
		return err
	}

	plist := plistPath()
	if err := os.MkdirAll(filepath.Dir(plist), 0o755); err != nil {
		return err
	}

	tmpl, err := template.New("plist").Parse(plistTemplate)
	if err != nil {
		return err
	}

	f, err := os.Create(plist)
	if err != nil {
		return err
	}
	defer f.Close()

	return tmpl.Execute(f, struct {
		BinaryPath string
		LogDir     string
	}{
		BinaryPath: binPath,
		LogDir:     logDir,
	})
}

// UninstallService removes the launchd plist.
func UninstallService() error {
	_ = exec.Command("launchctl", "unload", plistPath()).Run()
	return os.Remove(plistPath())
}
