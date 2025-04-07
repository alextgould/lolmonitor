// Run using: go test -v -count=1 ./internal/notifications
// (this will disable caching since this is a live test)

package notifications

import (
	"testing"
	"time"
)

func TestNotifications(t *testing.T) {
	SendNotification("test notification", "this is a test of the notification system", false)
	time.Sleep(5 * time.Second) // pause to give the notification time to appear
}

func TestDelayClose(t *testing.T) {
	delay := 20
	sessionGames := 2
	gamesPerSession := 3
	DelayClose(delay, sessionGames, gamesPerSession)
	time.Sleep(5 * time.Second) // pause to give the notification time to appear
}

func TestEndOfGame(t *testing.T) {
	endOfBreak := time.Now().Add(30 * time.Minute)
	sessionGames := 1
	gamesPerSession := 3
	EndOfGame(endOfBreak, sessionGames, gamesPerSession)
	time.Sleep(5 * time.Second) // pause to give the notification time to appear
}

func TestLobbyBlocked(t *testing.T) {
	endOfBreak := time.Now().Add(15 * time.Minute)
	LobbyBlocked(endOfBreak)
	time.Sleep(5 * time.Second) // pause to give the notification time to appear
}
