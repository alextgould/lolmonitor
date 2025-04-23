package utils

import (
	"os"
	"path/filepath"
)

// GetCurrentPath returns the directory of the currently running executable.
// If os.Executable() fails (e.g., during go run or tests), it falls back to the current working directory.
func GetCurrentPath(includeFilename bool) (string, error) {
	exePath, err := os.Executable()
	if err != nil {
		exePath, err = os.Getwd()
		if err != nil {
			return "", err
		}
	}

	exePath, err = filepath.Abs(exePath)
	if err != nil {
		return "", err
	}

	if includeFilename {
		return exePath, nil
	}
	return filepath.Dir(exePath), nil
}
