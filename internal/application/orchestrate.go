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

	var gameStartTime, gameEndTime, endOfBreak time.Time
	var breakDuration time.Duration
	var sessionGames int
	lastChecked := time.Now()

	for event := range events {

		// reset session games if there was a significant gap since the last game ended (and we're either opening the lobby or starting a new game)
		if sessionGames > 0 && event.Type == "open" && gameEndTime.Add(time.Duration(cfg.BreakBetweenSessionsMinutes)*time.Minute).Before(time.Now()) {
			sessionGames = 0
		}

		// update the config if it has been modified recently
		configUpdated, err := config.CheckConfigUpdated("", lastChecked)
		if err != nil {
			log.Printf("Error checking if config was updated: %v", err)
		} else if configUpdated {
			cfg, err = config.LoadConfig("")
			if err != nil {
				log.Panicf("Failed to load config: %v", err)
			}
		}
		lastChecked = time.Now()

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

					// optional delay until close takes effect
					if cfg.LobbyCloseDelaySeconds > 0 {
						notifications.DelayClose(cfg.LobbyCloseDelaySeconds)
						time.Sleep(time.Duration(cfg.LobbyCloseDelaySeconds) * time.Second)
					}

					log.Printf("Closing the lobby. Break will expire at: %v", endOfBreak.Format("2006/01/02 15:04:05"))
					err := window.Close(LOBBY_WINDOW_NAME, pr, pk)
					if err != nil {
						log.Printf("Error: %v", err)
					}
					notifications.EndOfGame(endOfBreak)
				}
			}
		} else if event.Type == "open" {
			log.Println("Lobby was opened")

			// check if Lobby is banned and force close if so
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

// TODO - need to adjust logic in orchestrate to reset the session games if a session break duration passes
// e.g. you play 2 games and can play up to 3 games. then you don't play.
// next day you come back and want to play your 3 games but it counts the 1st as your 3rd in the session
// that would be bad
