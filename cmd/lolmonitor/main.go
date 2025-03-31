// To run (in terminal, at project root): go run cmd/lolmonitor/main.go
// To build: go build cmd/lolmonitor/main.go
package main

import (
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/alextgould/lolmonitor/internal/actions"
	"github.com/alextgould/lolmonitor/internal/config"
	"github.com/alextgould/lolmonitor/internal/detection"
)

// import "github.com/go-toast/toast"

// func sendNotification(title, message string) {
// 	notification := toast.Notification{
// 		AppID:   "LoL Monitor",
// 		Title:   title,
// 		Message: message,
// 	}
// 	err := notification.Push()
// 	if err != nil {
// 		log.Printf("Failed to send notification: %v", err)
// 	}
// }

// sendNotification("League of Legends Blocked", "You're on a break! Try again later.")

func main() {
	// Load configuration
	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		log.Panicf("Failed to load config: %v", err)
	}

	// Start window monitoring
	log.Println("Starting lolmonitor")
	gameEvents := make(chan detection.GameEvent, 10) // Buffered channel to prevent blocking
	quit := make(chan struct{})
	go detection.StartMonitoring(gameEvents, quit)

	// Handle Ctrl+C for graceful shutdown
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		<-signalChan
		log.Println("Shutting down lolmonitor")
		close(quit) // Signal monitoring to stop
		os.Exit(0)
	}()

	var gameStartTime time.Time
	var gameEndTime time.Time
	var gameDuration time.Duration
	var breakDuration time.Duration
	var endOfBreak time.Time
	var sessionGames int

	for event := range gameEvents { // will process events as they are added to the gameEvents channel by the StartMonitoring goroutine, which runs constantly
		switch event.Type { // type of event (lobby_open, lobby_close, game_open, game_close)
		case "game_open": // note the game start time, so we can check for remakes
			log.Println("Game started")
			gameStartTime = event.Timestamp

		case "game_close": // check for remakes, increment game counter
			log.Println("Game ended")
			gameEndTime = event.Timestamp
			gameDuration = gameEndTime.Sub(gameStartTime)
			if gameDuration < time.Duration(cfg.MinimumGameDurationMinutes)*time.Minute {
				log.Println("Game was too short, likely a remake.")
				continue
			}
			sessionGames++
			log.Printf("Games played this session: %d of %d", sessionGames, cfg.GamesPerSession)
			if sessionGames >= cfg.GamesPerSession && cfg.GamesPerSession != 0 { // use 0 to disable GamesPerSession functionality
				sessionGames = 0
				breakDuration = time.Duration(cfg.BreakBetweenSessionsMinutes) * time.Minute
				log.Printf("Enforcing a session ban of %v minutes", cfg.BreakBetweenSessionsMinutes)
			} else {
				breakDuration = time.Duration(cfg.BreakBetweenGamesMinutes) * time.Minute
				log.Printf("Enforcing a game ban of %v minutes", cfg.BreakBetweenGamesMinutes)
			}
			endOfBreak = gameEndTime.Add(breakDuration)

			// always close the lobby after a game, unless breakDuration is 0
			if breakDuration > 0 {
				// TODO: provide a message to the user, using Windows toast system
				log.Printf("Closing the lobby. Break will expire at: %v", endOfBreak.Format("2006/01/02 15:04:05"))
				err := actions.CloseLeagueClient()
				if err != nil {
					log.Printf("Error: %v", err)
				}
			}

		case "lobby_open": // check if we are in a break, if so, close the lobby and inform the user
			log.Println("Lobby was opened")
			if time.Now().Before(endOfBreak) {
				log.Println("Closing the lobby, which was opened during an enforced break.")
				// TODO: provide a message to the user, using Windows toast system
				err := actions.CloseLeagueClient()
				if err != nil {
					log.Printf("Error: %v", err)
				}
			}

		case "lobby_close": // will happen when we close the lobby, or when the user closes the lobby, either way ignore
			continue

		default: // should never occur
			log.Panicf("Unknown event type: %v", event.Type)
		}
	}
}
