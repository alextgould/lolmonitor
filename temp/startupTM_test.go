// these tests require elevated privileges. hence, we can run
// go test -c -o ./internal/infrastructure/startup/startup_test.exe ./internal/infrastructure/startup
// to create a test executable binary
// then right click > run as administrator in windows explorer
package startupTM

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"
)

const TEST_TASK_NAME = "test_lolmonitor_task"

// Redirect log output to file
func TestMain(m *testing.M) {
	// Create a log file
	logFile, err := os.OpenFile("startup_test.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644) // create file if it doesn't exist, open for write only, append data instead of overwriting; 0644 = owner read and write, group read-only, others read-only

	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer logFile.Close()

	// Redirect log output to the file
	log.SetOutput(logFile)
	log.Println("Starting tests...")

	// Run tests
	m.Run()

	// Pause before exiting
	fmt.Println("Tests completed. Press Enter to exit...")
	fmt.Scanln()
}

func TestIsAdmin(t *testing.T) {
	err := isAdmin()
	if err != nil {
		log.Printf("Admin check failed: %v", err)
		t.Skip("Skipping tests because admin access is required.")
	} else {
		log.Println("Admin check passed.")
	}
}

// func TestIsAdmin(t *testing.T) {
// 	err := isAdmin()
// 	if err != nil {
// 		log.Printf("Need admin access to run Task Scheduler related tests: %v", err)
// 		t.Fail()
// 		//t.Errorf("Need admin access to run Task Scheduler related tests: %v", err)
// 	}
// }

// Check get user name is working
// func TestUser(t *testing.T) {
// 	user := os.Getenv("USERNAME")
// 	//t.Errorf("user is %v", user)
// 	//log.Printf("user is %v", user)
// 	//t.Fail()
// }

// Ensure task doesn't exist
func TestNoExists(t *testing.T) {

	// Ensure the task does not exist initially
	exists, err := taskExists(TEST_TASK_NAME)
	if err != nil {
		log.Printf("Error checking if task exists: %v", err)
		t.Fail()
	}
	if exists {
		log.Printf("Task %s should not exist initially", TEST_TASK_NAME)
		t.Fail()
	}
	log.Println("Completed taskExists check")
}

// Create the task
func TestCreate(t *testing.T) {
	// The error code 0x80004005 is a generic Windows error code that translates to "Unspecified error."
	exePath, err := getCurrentPath()
	if err != nil {
		log.Printf("Error getting current path: %v", err)
		//t.Errorf("Error getting current path: %v", err)
	}
	err = addTask(TEST_TASK_NAME, exePath)
	if err != nil {
		log.Printf("Error adding task: %v", err)
		//t.Errorf("Error adding task: %v", err)
	}
	log.Println("Completed addTask check")
	t.Fail() // check log outputs
}

// Verify the task now exists
func TestExists(t *testing.T) {
	exists, err := taskExists(TEST_TASK_NAME)
	if err != nil {
		t.Errorf("Error checking if task exists after creation: %v", err)
	}
	if !exists {
		t.Errorf("Task %s should exist after creation", TEST_TASK_NAME)
	}
}

// Remove the task
func TestRemove(t *testing.T) {
	err := removeTask(TEST_TASK_NAME)
	if err != nil {
		t.Errorf("Error removing task: %v", err)
	}
}

// Verify the task no longer exists
func TestNoExistsAfterRemove(t *testing.T) {
	exists, err := taskExists(TEST_TASK_NAME)
	if err != nil {
		t.Errorf("Error checking if task exists after removal: %v", err)
	}
	if exists {
		t.Errorf("Task %s should not exist after removal", TEST_TASK_NAME)
	}
}

func TestGetCurrentPath(t *testing.T) {
	// Get the current path
	exePath, err := getCurrentPath()
	if err != nil {
		t.Errorf("Error getting current path: %v", err)
	}

	// Verify the path matches the absolute path of this test file
	expectedPath, err := filepath.Abs(os.Args[0])
	if err != nil {
		t.Errorf("Error getting absolute path of test file: %v", err)
	}
	if exePath != expectedPath {
		t.Errorf("Expected path %s, but got %s", expectedPath, exePath)
	}
}
