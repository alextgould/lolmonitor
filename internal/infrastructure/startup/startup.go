// manage windows task scheduler tasks
package startup

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"golang.org/x/sys/windows/registry"

	"github.com/alextgould/lolmonitor/internal/interfaces/notifications"
)

// const but want to change it in _test.go
var TASK_NAME = "lolmonitor"

// get the current path of the executable
func getCurrentPath() (string, error) {
	exePath, err := os.Executable()
	if err != nil {
		return "", err
	}
	exePath, err = filepath.Abs(exePath)
	if err != nil {
		return "", err
	}
	return exePath, nil
}

// create registry key
func addToStartup(taskName, exePath string) error {
	key, _, err := registry.CreateKey(
		registry.CURRENT_USER,
		`Software\Microsoft\Windows\CurrentVersion\Run`,
		registry.SET_VALUE,
	)
	if err != nil {
		return fmt.Errorf("failed to open registry key: %v", err)
	}
	defer key.Close()

	err = key.SetStringValue(taskName, exePath)
	if err != nil {
		return fmt.Errorf("failed to set registry value: %v", err)
	}

	slog.Info("Successfully added to startup", "taskName", taskName, "exePath", exePath)
	return nil
}

// remove registry key
func removeFromStartup(taskName string) error {
	key, err := registry.OpenKey(
		registry.CURRENT_USER,
		`Software\Microsoft\Windows\CurrentVersion\Run`,
		registry.SET_VALUE,
	)
	if err != nil {
		return err
	}
	defer key.Close()

	return key.DeleteValue(taskName)
}

// check if regkey exists, returns a bool and the exePath if it exists
func startupEntryExists(taskName string) (bool, string, error) {
	key, err := registry.OpenKey(
		registry.CURRENT_USER,
		`Software\Microsoft\Windows\CurrentVersion\Run`,
		registry.QUERY_VALUE,
	)
	if err != nil {
		return false, "", err
	}
	defer key.Close()

	exePath, _, err := key.GetStringValue(taskName)
	if err == registry.ErrNotExist { // check worked and regkey doesn't exist
		return false, "", nil
	}
	if err != nil { // check failed
		return false, "", err
	}
	return true, exePath, nil // check worked and regkey exists
}

// if LoadOnStartup is True, confirm a valid registry entry exists or else add it
func ConfirmLoadOnStartup() error {

	taskExists, taskPath, err := startupEntryExists(TASK_NAME)
	if err != nil {
		slog.Error("an error occurred while checking if the startup registry entry exists")
		return err
	}

	// get location of .exe
	exePath, err := getCurrentPath()
	if err != nil {
		slog.Error("an error occurred while getting the current exe path")
		return err
	}

	// registry entry exists and path matches
	if taskExists {
		if taskPath == exePath {
			slog.Info("confirmed startup registry entry exists and is valid")
			return nil
		} else {
			slog.Info("startup registry entry exists but path is different - recreating registry entry")

			// remove the task
			err = removeFromStartup(TASK_NAME)
			if err != nil {
				slog.Error("error removing startup registry entry", "err", err)
				return err
			}
		}
	}

	// add task
	err = addToStartup(TASK_NAME, exePath)
	if err != nil {
		slog.Error("an error occured while adding startup registry entry", "err", err)
		return err
	}
	slog.Info("registry entry was added")
	notifications.SendNotification("Automatic startup added", "You no longer need to manually run the lolmonitor.exe file, it will load automatically. If you want to remove this, adjust the config file to have LoadOnStartup: false", false)
	return nil
}

// main function
func ConfirmNoLoadOnStartup() error {

	taskExists, _, err := startupEntryExists(TASK_NAME)
	if err != nil {
		slog.Error("an error occurred while checking if the startup registry entry exists")
		return err
	}

	if !taskExists {
		slog.Info("confirmed startup registry entry does not exist")
		return nil
	}

	// remove the task
	err = removeFromStartup(TASK_NAME)
	if err != nil {
		slog.Error("error removing startup registry entry", "err", err)
		return err
	}

	slog.Info("registry entry was removed")
	notifications.SendNotification("Automatic startup removed", "You now need to manually run lolmonitor.exe each time you play. If you want it to run automatically, adjust the config file to have LoadOnStartup: true", false)
	return nil
}
