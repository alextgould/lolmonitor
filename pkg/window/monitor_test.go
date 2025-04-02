package window

import (
	"testing"
	"time"
)

func TestIsProcessActive(t *testing.T) {
	mockRetriever := &MockProcessRetriever{Processes: make(map[string]Win32_Process)}

	// Add a fake process
	mockRetriever.AddProcess("My window", 1234)

	// Call monitoring function
	exists, PID, _ := isProcessActive("My window", mockRetriever)
	if !exists {
		t.Errorf("Expected window to be detected, but it wasn't")
	}
	if PID != 1234 {
		t.Errorf("Expected process ID of 1234 but got %v", PID)
	}

	// Remove the process
	mockRetriever.RemoveProcess("My window")

	// Call monitoring function again
	exists, _, _ = isProcessActive("My window", mockRetriever)
	if exists {
		t.Errorf("Expected window to be gone, but it was detected")
	}
}

func TestMonitorProcess(t *testing.T) {
	// A buffer of 10 allows up to 10 events to be queued without blocking the sender.
	// This is useful if multiple events are generated in quick succession.
	myEvents := make(chan ProcessEvent, 10)
	mockRetriever := &MockProcessRetriever{Processes: make(map[string]Win32_Process)}

	// Start monitoring "My window"
	go MonitorProcess("My window", myEvents, 1, mockRetriever) // will check processes every 1 second

	// Add a process and check for an "open" event
	mockRetriever.AddProcess("My window", 1234)
	time.Sleep(2 * time.Second) // Allow time for process to be detected and event sent
	select {
	case event := <-myEvents:
		if event.Type != "open" || event.Name != "My window" || event.PID != 1234 {
			t.Errorf("Expected 'open' event for 'My window' with PID 1234, got %+v", event)
		}
	case <-time.After(1 * time.Second):
		t.Errorf("Timed out waiting for 'open' event for 'My window'")
	}

	// Test running multiple MonitorProcess instances
	go MonitorProcess("My other window", myEvents, 1, mockRetriever)
	mockRetriever.AddProcess("My other window", 5678)
	time.Sleep(2 * time.Second) // Allow time for process to be detected and event sent
	select {
	case event := <-myEvents:
		if event.Type != "open" || event.Name != "My other window" || event.PID != 5678 {
			t.Errorf("Expected 'open' event for 'My other window' with PID 5678, got %+v", event)
		}
	case <-time.After(1 * time.Second):
		t.Errorf("Timed out waiting for 'open' event for 'My other window'")
	}

	// Remove a process and check for a "close" event
	mockRetriever.RemoveProcess("My window")
	time.Sleep(2 * time.Second) // Allow time for process to be detected and event sent
	select {
	case event := <-myEvents:
		if event.Type != "close" || event.Name != "My window" {
			t.Errorf("Expected 'close' event for 'My window', got %+v", event)
		}
	case <-time.After(1 * time.Second):
		t.Errorf("Timed out waiting for 'close' event for 'My window'")
	}

	// Close the events channel and wait for monitoring to end
	close(myEvents)
	time.Sleep(2 * time.Second) // Allow time for goroutines to exit

	// No additional confirmation is needed here as the goroutines will naturally terminate
	// when the channel is closed and no further events can be sent.
}

func TestWaitForProcessClose(t *testing.T) {
	mockRetriever := &MockProcessRetriever{Processes: make(map[string]Win32_Process)}

	// Add a process to simulate its existence
	mockRetriever.AddProcess("My window", 1234)

	// Use a goroutine to simulate process closure after 2 seconds
	go func() {
		time.Sleep(2 * time.Second)
		mockRetriever.RemoveProcess("My window")
	}()

	// Call WaitForProcessClose and verify it returns after the process is removed
	start := time.Now()
	WaitForProcessClose("My window", 1, mockRetriever)
	elapsed := time.Since(start)

	// Ensure the function waited approximately 2 seconds
	if elapsed < 2*time.Second || elapsed > 4*time.Second {
		t.Errorf("WaitForProcessClose did not wait the expected duration, elapsed: %v", elapsed)
	}
}
