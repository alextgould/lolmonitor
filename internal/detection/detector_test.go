// Run using: go test -v -timeout 2m -count=1 ./internal/detection
// -timeout 2m extends the test beyond the typical 30s limit (using vscode right click > "run tests")
//   so we have sufficient time to open/close the Lobby and/or practice Game
// -count=1 flag prevents Go from using cached test results, ensuring the test runs fresh every time

package detection

import (
	"testing"
	"time"
)

func TestLiveMonitoring(t *testing.T) {

	gameEvents := make(chan GameEvent, 10) // Buffered channel to prevent blocking
	quit := make(chan struct{})

	go StartMonitoring(gameEvents, quit)

	t.Log("Live monitoring started. Open and close the relevant windows to test.")

	// Ensure cleanup happens even if the test fails
	t.Cleanup(func() {
		close(quit)
		t.Log("Stopped monitoring.")
	})

	timeoutDuration := 90 * time.Second
	timeout := time.After(timeoutDuration)
	t.Logf("The test will run for %v seconds. Events will be logged below:", timeoutDuration.Seconds())
	for {
		select {
		case event := <-gameEvents:
			t.Logf("Detected event: %s at %v", event.Type, event.Timestamp)
		case <-timeout:
			t.Log("Live monitoring test completed.")
			return
		}
	}
}
