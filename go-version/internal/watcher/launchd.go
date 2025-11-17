package watcher

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"text/template"
)

const launchAgentLabel = "com.qkflow.watch"

const plistTemplate = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>{{.Label}}</string>
    
    <key>ProgramArguments</key>
    <array>
        <string>{{.ExecutablePath}}</string>
        <string>watch</string>
        <string>daemon</string>
    </array>
    
    <key>RunAtLoad</key>
    <true/>
    
    <key>KeepAlive</key>
    <dict>
        <key>SuccessfulExit</key>
        <false/>
    </dict>
    
    <key>StandardOutPath</key>
    <string>{{.StdoutPath}}</string>
    
    <key>StandardErrorPath</key>
    <string>{{.StderrPath}}</string>
    
    <key>WorkingDirectory</key>
    <string>{{.WorkingDir}}</string>
    
    <key>EnvironmentVariables</key>
    <dict>
        <key>PATH</key>
        <string>/usr/local/bin:/usr/bin:/bin:/usr/sbin:/sbin</string>
    </dict>
</dict>
</plist>
`

// LaunchAgentConfig holds configuration for launchd plist
type LaunchAgentConfig struct {
	Label          string
	ExecutablePath string
	StdoutPath     string
	StderrPath     string
	WorkingDir     string
}

// GetLaunchAgentPath returns the path to the launch agent plist file
func GetLaunchAgentPath() (string, error) {
	if runtime.GOOS != "darwin" {
		return "", fmt.Errorf("launchd is only supported on macOS")
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	launchAgentsDir := filepath.Join(homeDir, "Library", "LaunchAgents")
	return filepath.Join(launchAgentsDir, launchAgentLabel+".plist"), nil
}

// IsLaunchAgentInstalled checks if the launch agent is installed
func IsLaunchAgentInstalled() (bool, error) {
	plistPath, err := GetLaunchAgentPath()
	if err != nil {
		return false, err
	}

	_, err = os.Stat(plistPath)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}

// InstallLaunchAgent installs the launch agent
func InstallLaunchAgent(executablePath string) error {
	if runtime.GOOS != "darwin" {
		return fmt.Errorf("launchd is only supported on macOS")
	}

	// Get paths
	plistPath, err := GetLaunchAgentPath()
	if err != nil {
		return err
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	// Create LaunchAgents directory if it doesn't exist
	launchAgentsDir := filepath.Dir(plistPath)
	if err := os.MkdirAll(launchAgentsDir, 0755); err != nil {
		return fmt.Errorf("failed to create LaunchAgents directory: %w", err)
	}

	// Prepare log paths
	logDir := filepath.Join(homeDir, ".qkflow")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("failed to create log directory: %w", err)
	}

	stdoutPath := filepath.Join(logDir, "watch.stdout.log")
	stderrPath := filepath.Join(logDir, "watch.stderr.log")

	// Create plist content
	config := LaunchAgentConfig{
		Label:          launchAgentLabel,
		ExecutablePath: executablePath,
		StdoutPath:     stdoutPath,
		StderrPath:     stderrPath,
		WorkingDir:     homeDir,
	}

	tmpl, err := template.New("plist").Parse(plistTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse plist template: %w", err)
	}

	// Create plist file
	file, err := os.Create(plistPath)
	if err != nil {
		return fmt.Errorf("failed to create plist file: %w", err)
	}
	defer file.Close()

	if err := tmpl.Execute(file, config); err != nil {
		return fmt.Errorf("failed to write plist content: %w", err)
	}

	// Load the launch agent
	cmd := exec.Command("launchctl", "load", plistPath)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to load launch agent: %w (output: %s)", err, string(output))
	}

	// Start the service
	cmd = exec.Command("launchctl", "start", launchAgentLabel)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to start launch agent: %w (output: %s)", err, string(output))
	}

	return nil
}

// UninstallLaunchAgent uninstalls the launch agent
func UninstallLaunchAgent() error {
	if runtime.GOOS != "darwin" {
		return fmt.Errorf("launchd is only supported on macOS")
	}

	plistPath, err := GetLaunchAgentPath()
	if err != nil {
		return err
	}

	// Check if installed
	installed, err := IsLaunchAgentInstalled()
	if err != nil {
		return err
	}

	if !installed {
		return nil // Already uninstalled
	}

	// Stop the service
	cmd := exec.Command("launchctl", "stop", launchAgentLabel)
	cmd.Run() // Ignore errors, service might not be running

	// Unload the launch agent
	cmd = exec.Command("launchctl", "unload", plistPath)
	if output, err := cmd.CombinedOutput(); err != nil {
		// Ignore "Could not find specified service" errors
		if len(output) > 0 && string(output) != "" {
			// Log but don't fail
		}
	}

	// Remove plist file
	if err := os.Remove(plistPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove plist file: %w", err)
	}

	return nil
}

