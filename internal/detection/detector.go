package detection

import (
	"log"
	"time"

	"github.com/StackExchange/wmi"
)

type GameEvent struct {
	Type      string
	Timestamp time.Time
}

type Win32_Process struct {
	Name string
}

func StartMonitoring(gameEvents chan<- GameEvent, quit <-chan struct{}) {
	defer close(gameEvents)
	log.Println("Starting WMI-based window monitoring...")

	go monitorWindowEvents(gameEvents)

	// Block until quit signal is received
	<-quit
	log.Println("Stopping monitoring...")
}

func monitorWindowEvents(gameEvents chan<- GameEvent) {
	activeWindows := make(map[string]bool)

	for {
		processes, err := getActiveProcesses()
		if err != nil {
			log.Printf("WMI query failed: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}

		// Store which League-related processes are active
		newActiveWindows := map[string]bool{}
		for _, process := range processes {
			newActiveWindows[process.Name] = true
		}

		// Flag new processes (useful for confirming window names)
		// for _, process := range processes {
		// 	if !activeWindows[process.Name] {
		// 		log.Printf("New process detected: %s", process.Name)
		// 	}
		// }

		// Detect lobby open
		if !activeWindows["LeagueClientUx.exe"] && newActiveWindows["LeagueClientUx.exe"] {
			gameEvents <- GameEvent{Type: "lobby_open", Timestamp: time.Now()}
			go func() { // Run in a goroutine otherwise waitForProcessClose will block game start detection
				waitForProcessClose("LeagueClientUx.exe")
				gameEvents <- GameEvent{Type: "lobby_close", Timestamp: time.Now()}
			}()
		}

		// Detect game start
		if !activeWindows["League of Legends.exe"] && newActiveWindows["League of Legends.exe"] {
			gameEvents <- GameEvent{Type: "game_start", Timestamp: time.Now()}
			go func() { // Run in a goroutine
				waitForProcessClose("League of Legends.exe")
				gameEvents <- GameEvent{Type: "game_end", Timestamp: time.Now()}
			}()
		}

		// Update activeWindows to reflect the new state
		activeWindows = newActiveWindows

		time.Sleep(1 * time.Second) // Short delay to avoid excessive CPU usage
	}
}

func getActiveProcesses() ([]Win32_Process, error) {
	var processes []Win32_Process
	query := "SELECT Name FROM Win32_Process" // WHERE Name = 'League of Legends (TM) Client.exe' OR Name = 'LeagueClientUx.exe'"
	err := wmi.Query(query, &processes)
	return processes, err
}

func waitForProcessClose(processName string) {
	// log.Printf("Waiting for %s to close...", processName)

	for {
		processes, err := getActiveProcesses()
		if err != nil {
			log.Printf("Error querying WMI: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}

		// Check if the process has exited
		found := false
		for _, process := range processes {
			if process.Name == processName {
				found = true
				break
			}
		}

		if !found {
			log.Printf("%s has closed.", processName)
			break
		}

		time.Sleep(1 * time.Second) // Avoid excessive polling
	}
}
