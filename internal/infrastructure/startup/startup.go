// manage windows task scheduler tasks
package startup

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"golang.org/x/sys/windows/registry"

	"github.com/alextgould/lolmonitor/internal/interfaces/notifications"
	//"github.com/alextgould/lolmonitor/internal/infrastructure/logger"
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

// check if regkey exists
func startupEntryExists(taskName string) (bool, error) {
	key, err := registry.OpenKey(
		registry.CURRENT_USER,
		`Software\Microsoft\Windows\CurrentVersion\Run`,
		registry.QUERY_VALUE,
	)
	if err != nil {
		return false, err
	}
	defer key.Close()

	_, _, err = key.GetStringValue(taskName)
	if err == registry.ErrNotExist {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// main function
func ConfirmStatus(LoadOnStartup, StartupInstalled bool) error {

	// we want to confirm there's no automatic startup and uninstall it if there is
	if !LoadOnStartup {

		// check if task exists, in case user removed the StartupInstalled config line
		taskExists, err := startupEntryExists(TASK_NAME)
		if err != nil {
			slog.Error("error checking if startup registry entry exists")
			return err
		}
		if !taskExists {
			slog.Info("confirmed no startup registry entry")
			return nil
		}

		// remove the task
		err = removeFromStartup(TASK_NAME)
		if err != nil {
			slog.Error("error removing startup registry entry", "err", err)
			return err
		}

		// close the program (assume it was started by the task scheduler)
		slog.Info("registry entry was removed. closing program.")
		notifications.SendNotification("Automatic startup removed", "You now need to manually run lolmonitor.exe each time you play. If you want it to run automatically, adjust the config file to have LoadOnStartup: true")
		// os.Exit(0)
	}

	// we want to add a registry entry so the process starts automatically in the future
	if LoadOnStartup && !StartupInstalled {

		// get location of .exe so task scheduler can run the program
		exePath, err := getCurrentPath()
		if err != nil {
			return err
		}

		// add task
		err = addToStartup(TASK_NAME, exePath)
		if err != nil {
			return err
		}
		slog.Info("registry entry was added")
		notifications.SendNotification("Automatic startup added", "You no longer need to manually run the lolmonitor.exe file, it will load automatically. If you want to remove this, adjust the config file to have LoadOnStartup: false")
	}
	return nil
}
