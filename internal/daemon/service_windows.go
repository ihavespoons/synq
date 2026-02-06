//go:build windows

package daemon

import "fmt"

// InstallService is a stub on Windows.
func InstallService(configDir string) error {
	return fmt.Errorf("automatic service installation is not supported on Windows; use Task Scheduler manually")
}

// UninstallService is a stub on Windows.
func UninstallService() error {
	return fmt.Errorf("automatic service removal is not supported on Windows")
}
