// manage windows task scheduler tasks
package startup

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/alextgould/lolmonitor/internal/interfaces/notifications"
)

const TASK_NAME = "lolmonitor"

// check to see if we have elevated privileges, otherwise we cannot adjust the task scheduler settings
// hence, to test this, we can run
// go test -c -o startup_test.exe
// to create a test executable binary then try right click > run as administrator in windows explorer
func isAdmin() error {
	cmd := exec.Command("net", "session")
	return cmd.Run()
}

// check to see if we need to install or uninstall lolmonitor from Windows Task Scheduler
func ConfirmStatus(LoadOnStartup, StartupInstalled bool) error {

	// uninstall check
	if !LoadOnStartup {

		// check if task exists, in case user removed the StartupInstalled config line
		taskExists, err := taskExists(TASK_NAME)
		if err != nil {
			log.Printf("Error checking if task exists in Task Scheduler: %v", err)
			return err
		}
		if !taskExists { // default outcome if LoadOnStartup is false (or missing))
			return nil
		}

		// remove the task
		err = removeTask(TASK_NAME)
		if err != nil {
			log.Printf("Error removing task from Task Scheduler: %v", err)
			return err
		}

		// close the program (assume it was started by the task scheduler)
		log.Println("Task was uninstalled. Closing program.")
		notifications.SendNotification("lolmonitor automatic startup removed", "You now need to manually run lolmonitor.exe each time you play. If you want it to run automatically, adjust the config file to have LoadOnStartup: true")
		os.Exit(0)
	}

	// install check
	if LoadOnStartup && !StartupInstalled {

		// get location of .exe so task scheduler can run the program
		exePath, err := getCurrentPath()
		if err != nil {
			return err
		}

		// add task
		err = addTask(TASK_NAME, exePath)
		if err != nil {
			return err
		}

		notifications.SendNotification("lolmonitor automatic startup added", "You no longer need to manually run the lolmonitor.exe file, it will load automatically. If you want to remove this, adjust the config file to have LoadOnStartup: false")
	}
	return nil
}

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

// add task
func addTask(taskName, exePath string) error {
	// schtasks /Create /SC ONLOGON by default creates the task in the system-wide task library,
	// not scoped to just the current user — so it will require admin privileges, even though
	// you’re targeting a user-level action. Hence, Add the /RU (run as user) flag with the
	// current user’s name to explicitly scope it to the user context
	// I feel like ChatGPT might be hallucinating...
	// /RU does not expect the USERNAME environment variable by itself — it expects the full
	// user principal name, i.e., "DOMAIN\\username" or "COMPUTERNAME\\username" for local users.

	//user := os.Getenv("USERNAME")
	//user := fmt.Sprintf("%s\\%s", os.Getenv("COMPUTERNAME"), os.Getenv("USERNAME"))
	cmd := exec.Command("schtasks", "/Create",
		"/SC", "ONLOGON",
		"/TN", taskName,
		"/TR", exePath,
		//"/RU", user,
	)
	// return cmd.Run()

	// Can capture stderr for debugging
	stderr, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("addTask failed. error: %v stderr: %s", err, string(stderr))
	}
	log.Printf("addTask succeeded. stdout/stderr: %s", string(stderr))
	return nil
}

// remove task
func removeTask(taskName string) error {
	user := os.Getenv("USERNAME")
	cmd := exec.Command("schtasks", "/Delete",
		"/TN", taskName,
		"/RU", user,
	)
	return cmd.Run()
}

// check if task exists
func taskExists(taskName string) (bool, error) {
	user := os.Getenv("USERNAME")
	cmd := exec.Command("schtasks", "/Query",
		"/TN", taskName,
		"/RU", user,
	)
	err := cmd.Run()

	if err != nil {
		// Check if the error is an ExitError with exit code 1 (task not found)
		if exitError, ok := err.(*exec.ExitError); ok && exitError.ExitCode() == 1 {
			return false, nil // Task does not exist
		}
		// Return other errors (e.g., command execution issues)
		return false, err
	}

	// Task exists
	return true, nil
}

// import (
// 	"golang.org/x/sys/windows/registry"
// )

// func addToStartup(taskName, exePath string) error {
// 	key, _, err := registry.CreateKey(
// 		registry.CURRENT_USER,
// 		`Software\Microsoft\Windows\CurrentVersion\Run`,
// 		registry.SET_VALUE,
// 	)
// 	if err != nil {
// 		return err
// 	}
// 	defer key.Close()

// 	return key.SetStringValue(taskName, exePath)
// }

// import (
// 	"golang.org/x/sys/windows/registry"
// )

// func removeFromStartup(taskName string) error {
// 	key, err := registry.OpenKey(
// 		registry.CURRENT_USER,
// 		`Software\Microsoft\Windows\CurrentVersion\Run`,
// 		registry.SET_VALUE,
// 	)
// 	if err != nil {
// 		return err
// 	}
// 	defer key.Close()

// 	return key.DeleteValue(taskName)
// }
