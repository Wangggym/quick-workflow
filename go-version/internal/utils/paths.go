package utils

import (
	"os"
	"path/filepath"
	"runtime"
)

// GetConfigDir returns the config directory, preferring iCloud Drive on macOS
func GetConfigDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	// On macOS, try to use iCloud Drive first
	if runtime.GOOS == "darwin" {
		iCloudPath := filepath.Join(homeDir, "Library", "Mobile Documents", "com~apple~CloudDocs", ".qkflow")
		
		// Check if iCloud Drive is available
		iCloudBase := filepath.Join(homeDir, "Library", "Mobile Documents", "com~apple~CloudDocs")
		if info, err := os.Stat(iCloudBase); err == nil && info.IsDir() {
			// iCloud Drive is available, create our config dir if needed
			if err := os.MkdirAll(iCloudPath, 0755); err == nil {
				return iCloudPath, nil
			}
		}
	}

	// Fallback to local directory
	localPath := filepath.Join(homeDir, ".qkflow")
	if err := os.MkdirAll(localPath, 0755); err != nil {
		return "", err
	}
	
	return localPath, nil
}

// GetQuickWorkflowConfigDir returns the quick-workflow config directory
// This is used for the main config.yaml file
// Note: This now returns the same directory as GetConfigDir() for consistency
func GetQuickWorkflowConfigDir() (string, error) {
	return GetConfigDir()
}

// IsICLoudAvailable checks if iCloud Drive is available on macOS
func IsICLoudAvailable() bool {
	if runtime.GOOS != "darwin" {
		return false
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return false
	}

	iCloudBase := filepath.Join(homeDir, "Library", "Mobile Documents", "com~apple~CloudDocs")
	info, err := os.Stat(iCloudBase)
	return err == nil && info.IsDir()
}

// GetConfigLocation returns a human-readable description of where configs are stored
func GetConfigLocation() string {
	if IsICLoudAvailable() {
		return "iCloud Drive (synced across devices)"
	}
	return "Local storage"
}

