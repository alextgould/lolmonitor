package application

import (
	"log"
	"time"

	"github.com/alextgould/lolmonitor/internal/config"
	"github.com/alextgould/lolmonitor/internal/interfaces/notifications"

	window "github.com/alextgould/lolmonitor/pkg/window"
)

const (
	LOBBY_WINDOW_NAME       = "LeagueClientUx.exe"
	GAME_WINDOW_NAME        = "League of Legends.exe"
	CHECK_FREQUENCY_SECONDS = 15
)

// monitor for Lobby and Game events, using the lobbyEvents and gameEvents channels
// use the config settings with domain logic to force close windows as needed
// pr and pk should be nil unless running tests
func Monitor(cfg config.Config, events chan window.ProcessEvent, pr window.ProcessRetriever, pk window.ProcessKiller) {

	// use WMI by default (unless running tests)
	if pr == nil {
		pr = window.WMIProcessRetriever{}
	}
	if pk == nil {
		pk = window.WMIProcessKiller{}
	}

	// monitor for the Lobby and/or Game windows opening and closing
	go window.MonitorProcess(LOBBY_WINDOW_NAME, events, CHECK_FREQUENCY_SECONDS, pr)
	log.Println("Monitoring for Lobby window")
	go window.MonitorProcess(GAME_WINDOW_NAME, events, CHECK_FREQUENCY_SECONDS, pr)
	log.Println("Monitoring for Game window")

	var gameStartTime time.Time
	var gameEndTime time.Time
	var breakDuration time.Duration
	var endOfBreak time.Time
	var sessionGames int

	for event := range events {
		log.Printf("Processing event - type: %v name: %v", event.Type, event.Name)
		if event.Name == GAME_WINDOW_NAME {
			if event.Type == "open" {
				log.Println("Game started")
				gameStartTime = event.Timestamp
				window.WaitForProcessClose(GAME_WINDOW_NAME, 1, pr) // check frequently
			} else {
				log.Println("Game ended")
				gameEndTime = event.Timestamp
				sessionGames, endOfBreak, breakDuration = postGame(cfg, time.Now(), gameStartTime, gameEndTime, sessionGames)

				// always close the lobby after a game, unless breakDuration is 0
				if breakDuration > 0 {
					log.Printf("Closing the lobby, game ended. Break will expire at: %v", endOfBreak.Format("2006/01/02 15:04:05"))

					err := window.Close(LOBBY_WINDOW_NAME, pr, pk)
					if err != nil {
						log.Printf("Error: %v", err)
					}
					notifications.EndOfGame(endOfBreak)
				}
			}
		} else if event.Type == "open" {
			log.Println("Lobby was opened")
			isLobbyBan, err := isLobbyBan(cfg, time.Now(), endOfBreak)
			if err != nil {
				log.Printf("Error: %v", err)
			} else if isLobbyBan {
				log.Println("Closing the lobby, which was opened during an enforced break.")
				err := window.Close(LOBBY_WINDOW_NAME, pr, pk)
				if err != nil {
					log.Printf("Error: %v", err)
				}
				notifications.LobbyBlocked(endOfBreak)
			}
		}
	}
}

// post game logic to determine the new sessionGames and endOfBreak values
// Confirming: Go does not allow implicit pointer updates for integers like sessionGames.
// You must return the updated value and assign it in the calling function.
func postGame(cfg config.Config, currentTime, gameStartTime, gameEndTime time.Time, sessionGames int) (int, time.Time, time.Duration) {
	var gameDuration, breakDuration time.Duration

	gameDuration = gameEndTime.Sub(gameStartTime)

	// increment game session count
	if gameDuration < time.Duration(cfg.MinimumGameDurationMinutes)*time.Minute {
		log.Println("Game was below the minimum duration (remake).")
		return sessionGames, currentTime, 0
	}
	sessionGames++

	// update required break duration based on config settings
	log.Printf("Games played this session: %d of %d", sessionGames, cfg.GamesPerSession)
	if sessionGames >= cfg.GamesPerSession && cfg.GamesPerSession != 0 { // use 0 to disable GamesPerSession functionality
		sessionGames = 0
		breakDuration = time.Duration(cfg.BreakBetweenSessionsMinutes) * time.Minute
		log.Printf("Enforcing a session ban of %v minutes", cfg.BreakBetweenSessionsMinutes)
	} else {
		breakDuration = time.Duration(cfg.BreakBetweenGamesMinutes) * time.Minute
		log.Printf("Enforcing a game ban of %v minutes", cfg.BreakBetweenGamesMinutes)
	}
	endOfBreak := gameEndTime.Add(breakDuration)
	return sessionGames, endOfBreak, breakDuration
}

// when the Lobby window is opened, check if it's allowed to be opened
func isLobbyBan(cfg config.Config, currentTime, endOfBreak time.Time) (bool, error) {

	// check we are not the endOfBreak exclusion period
	if currentTime.Before(endOfBreak) {
		log.Println("current time before endofbreak, return true ban")
		return true, nil
	}

	// check if we are not outside of the [DailyStartTime, DailyEndTime] allowed period
	currentHM := currentTime.Format("15:04")
	if (cfg.DailyStartTime != "00:00" && cfg.DailyStartTime != "" && currentHM <= cfg.DailyStartTime) || (cfg.DailyEndTime != "00:00" && cfg.DailyEndTime != "" && currentHM >= cfg.DailyEndTime) {
		log.Println("outside of the DailyStartTime / DailyEndTime range")
		return true, nil
	}

	// otherwise the Lobby is not banned
	log.Println("lobby is not banned")
	return false, nil
}
