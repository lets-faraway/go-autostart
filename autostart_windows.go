package autostart

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const powershellCreateShortcutTemplate = `$WshShell = New-Object -comObject WScript.Shell
$Shortcut = $WshShell.CreateShortcut("%s")
$Shortcut.TargetPath = "%s"
$Shortcut.Arguments = "%s"
$Shortcut.Save()`

var startupDir string

func init() {
	startupDir = filepath.Join(os.Getenv("USERPROFILE"), "AppData", "Roaming", "Microsoft", "Windows", "Start Menu", "Programs", "Startup")
}

func (a *App) path() string {
	return filepath.Join(startupDir, a.Name+".lnk")
}

func (a *App) createShortcut(lnkPath string, targetPath string, arg string) error {
	// PowerShell script to create a shortcut
	psScript := fmt.Sprintf(powershellCreateShortcutTemplate, lnkPath, targetPath, arg)

	// Execute the PowerShell script
	cmd := exec.Command("powershell", "-Command", psScript)
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func (a *App) IsEnabled() bool {
	_, err := os.Stat(a.path())
	return err == nil
}

func (a *App) Enable() error {
	if err := os.MkdirAll(startupDir, 0777); err != nil {
		return err
	}
	if err := a.createShortcut(a.path(), a.Exec[0], strings.Join(a.Exec[1:], " ")); err != nil {
		return errors.New(fmt.Sprintf("autostart: cannot create shortcut '%s'", a.path()))
	}
	return nil
}

func (a *App) Disable() error {
	return os.Remove(a.path())
}
