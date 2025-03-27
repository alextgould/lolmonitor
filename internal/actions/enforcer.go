package actions

import (
	"fmt"
	"lolmonitor/internal/config"
	"time"
)

// EnforceBreak closes the League of Legends client and notifies the user
func EnforceBreak(cfg config.Config, lastGameEnd time.Time, breakType string) {
	var breakDuration time.Duration

	if breakType == "game" {
		breakDuration = cfg.BreakBetweenGames
	} else if breakType == "session" {
		breakDuration = cfg.BreakBetweenSessions
	} else {
		fmt.Println("Unknown break type")
		return
	}

	endOfBreak := lastGameEnd.Add(breakDuration)
	fmt.Printf("Enforcing %s break until %v. Closing League client...\n", breakType, endOfBreak)

	// Call function to close the League Client window
	closeLeagueClient()
}

// closeLeagueClient forcefully closes the League client
func closeLeagueClient() {
	// TODO: Implement logic to find and close the League Client process
	fmt.Println("League of Legends client closed.")
}
