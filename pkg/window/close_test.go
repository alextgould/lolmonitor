package window

import (
	"testing"
)

func TestClose(t *testing.T) {
	pr := &MockProcessRetriever{Processes: make(map[string]Win32_Process)}
	pk := &MockProcessKiller{}

	// Add a fake process
	pr.AddProcess("My window", 1234)

	// Kill the process
	err := Close("My window", pr, pk)
	if err != nil {
		t.Fatalf("Unable to close the window: %v", err)
	}
	pr.RemoveProcess("My window")

	// Check process was removed
	exists, _, _ := isProcessActive("My window", pr)
	if exists {
		t.Errorf("Expected window to be gone, but it was detected")
	}
}
