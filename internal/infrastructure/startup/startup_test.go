// these tests require elevated privileges. hence, we can run
// go test -c -o ./internal/infrastructure/startup/startup_test.exe ./internal/infrastructure/startup
// to create a test executable binary
// then right click > run as administrator in windows explorer
package startup

import (
	"log/slog"
	"os"
	"path/filepath"
	"testing"
	"time"
)

const TEST_TASK_NAME = "test_lolmonitor_task"

// Show DEBUG log messages
func TestMain(m *testing.M) {

	// show debug log messages
	testLogger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(testLogger)

	// Run tests
	m.Run()
}

// test internal functions

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

// note you can manually run these tests out of order
// to gain confidence the functions actually work
// (e.g. run TestAddToStartup > TestStartupEntryNoExists to make it fail)
// you can also browse to regedit > Computer\HKEY_CURRENT_USER\SOFTWARE\Microsoft\Windows\CurrentVersion\Run
// and observe the task being created and removed as you run each test manually

func TestAddToStartup(t *testing.T) {
	exePath, err := getCurrentPath()
	if err != nil {
		t.Errorf("Error getting current path: %v", err)
		return

	}

	err = addToStartup(TEST_TASK_NAME, exePath)
	if err != nil {
		t.Errorf("Error adding task: %v", err)
		return
	}
}

func TestStartupEntryExists(t *testing.T) {
	exists, _, err := startupEntryExists(TEST_TASK_NAME)
	if err != nil {
		t.Errorf("startupEntryExists error: %v", err)
	}
	if !exists {
		t.Errorf("startup entry %s should exist after creation", TEST_TASK_NAME)
	}
	t.Log("Finished the startupEntryExists test")
}

func TestRemoveFromStartup(t *testing.T) {
	err := removeFromStartup(TEST_TASK_NAME)
	if err != nil {
		t.Errorf("removeFromStartup error: %v", err)
	}
	t.Log("Finished the RemoveFromStartup test")
}

func TestStartupEntryNoExists(t *testing.T) {
	exists, _, err := startupEntryExists(TEST_TASK_NAME)
	if err != nil {
		t.Errorf("startupEntryNoExists error: %v", err)
	}
	if exists {
		t.Errorf("startup entry %s should no longer exist after creation", TEST_TASK_NAME)
	}
	slog.Debug("Finished the startupEntryNoExists test")
}

// test external function

func TestConfirmFunctions(t *testing.T) {

	// override
	TASK_NAME = TEST_TASK_NAME

	// this should install it
	ConfirmLoadOnStartup()

	// check installed
	TestStartupEntryExists(t)

	// short break so I can manually see the process appear in regedit.exe
	time.Sleep(10 * time.Second)

	// this should do nothing
	ConfirmLoadOnStartup()

	// check installed
	TestStartupEntryExists(t)

	// this should uninstall it
	ConfirmNoLoadOnStartup()

	// check not installed
	TestStartupEntryNoExists(t)

	// this should do nothing
	ConfirmNoLoadOnStartup()

	// check not installed
	TestStartupEntryNoExists(t)
}
