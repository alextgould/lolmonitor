// Run using: go test -v -count=1 ./internal/notifications
// (this will disable caching since this is a live test)

package notifications

import (
	"testing"
	"time"
)

func TestNotifications(t *testing.T) {
	SendNotification("test notification", "this is a test of the notification system")
	time.Sleep(5 * time.Second) // pause to give the notification time to appear
}

func TestEndOfGame(t *testing.T) {
	endOfBreak := time.Now().Add(30 * time.Minute)
	EndOfGame(endOfBreak)
	time.Sleep(5 * time.Second) // pause to give the notification time to appear
}

func TestLobbyBlocked(t *testing.T) {
	endOfBreak := time.Now().Add(15 * time.Minute)
	LobbyBlocked(endOfBreak)
	time.Sleep(5 * time.Second) // pause to give the notification time to appear
}
