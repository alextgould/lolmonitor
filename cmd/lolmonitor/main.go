// cmd/lolmonitor/main.go
package main

import (
	"fmt"
	"log"
	"lolmonitor/internal/actions"
	"lolmonitor/internal/config"
	"lolmonitor/internal/detection"
	"time"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Start window monitoring
	fmt.Println("Starting League of Legends monitor...")
	gameEvents := detection.StartMonitoring()

	var lastGameEnd time.Time
	var sessionGames int

	for event := range gameEvents {
		if event.Type == "game_start" {
			fmt.Println("Game started at", event.Timestamp)
		} else if event.Type == "game_end" {
			fmt.Println("Game ended at", event.Timestamp)
			gameDuration := event.Timestamp.Sub(lastGameEnd)
			if gameDuration < cfg.MinGameDuration {
				fmt.Println("Game was too short, likely a remake. Ignoring.")
				continue
			}
			lastGameEnd = event.Timestamp
			sessionGames++
			if sessionGames >= cfg.GamesPerSession {
				actions.EnforceBreak(cfg.BreakBetweenSessions, "session")
				sessionGames = 0
			} else {
				actions.EnforceBreak(cfg.BreakBetweenGames, "game")
			}
		}
	}
}
