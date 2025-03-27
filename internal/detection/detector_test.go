// Run using: go test -v -count=1 ./internal/detection
// (this will disable caching since this is a live test)

// Confirmed working test results:
// === RUN   TestLiveMonitoring
//     detector_test.go:17: Live monitoring started. Open and close the relevant windows to test.
//     detector_test.go:18: The test will run for 1 minute. Events will be logged below:
// 2025/03/27 22:29:52 Starting WMI-based window monitoring...
//     detector_test.go:30: Detected event: lobby_open at 2025-03-27 22:29:52.4723685 +1100 AEDT m=+0.157270001
//     detector_test.go:30: Detected event: game_start at 2025-03-27 22:30:17.096469 +1100 AEDT m=+24.781370501
// 2025/03/27 22:30:28 League of Legends.exe has closed.
//     detector_test.go:30: Detected event: game_end at 2025-03-27 22:30:28.9147827 +1100 AEDT m=+36.599684201
// 2025/03/27 22:30:34 LeagueClientUx.exe has closed.
//     detector_test.go:30: Detected event: lobby_close at 2025-03-27 22:30:34.2980974 +1100 AEDT m=+41.982998901
//     detector_test.go:32: Live monitoring test completed.
//     detector_test.go:23: Stopped monitoring.
// 2025/03/27 22:30:52 Stopping monitoring...
// --- PASS: TestLiveMonitoring (60.00s)

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
	t.Log("The test will run for 1 minute. Events will be logged below:")

	// Ensure cleanup happens even if the test fails
	t.Cleanup(func() {
		close(quit)
		t.Log("Stopped monitoring.")
	})

	timeout := time.After(1 * time.Minute)
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
