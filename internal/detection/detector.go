package detection

import (
	"fmt"
	"log"
	"time"

	"github.com/StackExchange/wmi"
)

const (
	LOBBY_WINDOW_NAME = "LeagueClientUx.exe"
	GAME_WINDOW_NAME  = "League of Legends.exe"
)

// WMI package returns a slice of processes
type Win32_Process struct {
	Name string
}

// This package logs type of event (lobby_open, lobby_close, game_open, game_close), along with when it occured and the relevant process ID
type GameEvent struct {
	Type      string
	Timestamp time.Time
	ProcessId uint32
}

// Start the monitoring cycle
func StartMonitoring(gameEvents chan<- GameEvent, quit <-chan struct{}) {
	defer close(gameEvents)
	log.Println("Starting WMI-based window monitoring...")

	go monitorWindowEvents(gameEvents)

	// Block until quit signal is received
	<-quit
	log.Println("Stopping monitoring...")
}

// Monitor for LoL related windows being opened or closed
func monitorWindowEvents(gameEvents chan<- GameEvent) {

	lobbyWasActive := false
	for {
		lobbyActive, gameActive, _ := activeWindows()

		if !lobbyWasActive { // we are monitoring for the Lobby being opened
			if lobbyActive {
				gameEvents <- GameEvent{Type: "lobby_open", Timestamp: time.Now()}
			}
		} else if gameActive { // we are monitoring for Game being opened
			gameEvents <- GameEvent{Type: "game_open", Timestamp: time.Now()}
			waitForGameToEnd()
			gameEvents <- GameEvent{Type: "game_close", Timestamp: time.Now()}
		} else if !lobbyActive { // we are monitoring for Lobby being closed
			gameEvents <- GameEvent{Type: "lobby_close", Timestamp: time.Now()}
		}
		time.Sleep(10 * time.Second)
		lobbyWasActive = lobbyActive
	}
}

// use WMI to identify active LoL related windows
// returns bools indicating if the Lobby and Game windows are active respectively, as well as their process IDs
func activeWindows() (bool, bool, error) {
	var processes []Win32_Process
	query := fmt.Sprintf("SELECT Name FROM Win32_Process WHERE Name = '%s' OR Name = '%s'", LOBBY_WINDOW_NAME, GAME_WINDOW_NAME)

	// Retry loop with 5-second delay in case of WMI errors
	for {
		err := wmi.Query(query, &processes)
		if err == nil {
			var isLobbyRunning, isGameRunning bool
			for _, process := range processes {
				if process.Name == LOBBY_WINDOW_NAME {
					isLobbyRunning = true
				}
				if process.Name == GAME_WINDOW_NAME {
					isGameRunning = true
				}
			}
			return isLobbyRunning, isGameRunning, nil // Return the bools and no error
		}
		log.Printf("WMI query failed: %v", err)
		time.Sleep(5 * time.Second)
	}
}

func waitForGameToEnd() {
	for {
		_, gameActive, _ := activeWindows()
		if !gameActive {
			break
		}
		time.Sleep(1 * time.Second)
	}
}
