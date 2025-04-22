package utils

import (
	"os"
	"path/filepath"
)

// GetCurrentPath returns the directory of the currently running executable.
// If os.Executable() fails (e.g., during go run or tests), it falls back to the current working directory.
func GetCurrentPath() (string, error) {
	exePath, err := os.Executable()
	if err != nil {
		return os.Getwd()
	}
	exePath, err = filepath.Abs(exePath)
	if err != nil {
		return "", err
	}
	return filepath.Dir(exePath), nil
}
