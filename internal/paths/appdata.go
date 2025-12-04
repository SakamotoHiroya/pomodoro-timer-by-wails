package paths

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

func AppDataPath(appName string) (string, error) {
	var base string

	switch runtime.GOOS {
	case "windows":
		// %AppData%
		base = os.Getenv("APPDATA")
		if base == "" {
			// fallback: C:/Users/<User>/AppData/Roaming
			home, _ := os.UserHomeDir()
			base = filepath.Join(home, "AppData", "Roaming")
		}

	case "darwin":
		// ~/Library/Application Support
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		base = filepath.Join(home, "Library", "Application Support")

	default:
		// Linux or others: ~/.config
		// または ~/.local/share にする案もOK
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		base = filepath.Join(home, ".config")
	}

	if appName != "" {
		base = filepath.Join(base, appName)
	}

	// ensure directory exists
	if err := os.MkdirAll(base, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create app data directory: %w", err)
	}

	return base, nil
}
