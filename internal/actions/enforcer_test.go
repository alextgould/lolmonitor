// Run using: go test -v -count=1 ./internal/actions
// (this will disable caching since this is a live test)

package actions

import (
	"testing"
	"time"
)

func TestCloseLeagueClient(t *testing.T) {

	// Step 1: Ensure LeagueClientUx.exe is running before testing
	_, err := getLobbyPID()
	if err != nil {
		t.Skip("League Client is not running. Start it manually before running the test.")
	}

	// Step 2: Wait for 10 seconds before closing the client
	t.Log("League client found. Waiting 5 seconds before terminating it.")
	time.Sleep(5 * time.Second)

	// Step 3: Attempt to close the League client
	err = CloseLeagueClient()
	if err != nil {
		t.Errorf("Failed to close League client: %v", err)
	}

	// Step 4: Verify that the process is no longer running
	t.Log("League client closed. Waiting 5 seconds before confirming.")
	time.Sleep(5 * time.Second) // Give time for the process to terminate
	_, err = getLobbyPID()
	if err == nil {
		t.Errorf("League client is still running.")
	} else {
		t.Log("League client was successfully closed.")
		return
	}
}
