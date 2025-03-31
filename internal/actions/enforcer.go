package actions

import (
	"fmt"
	"log"
	"time"

	"github.com/alextgould/lolmonitor/internal/detection"

	"github.com/StackExchange/wmi"

	"os/exec"
)

func CloseLeagueClient() error {
	// var PID uint32
	PID, err := getLobbyPID()
	if err != nil {
		return fmt.Errorf("failed to close League client: %v", err)
	}
	cmd := exec.Command("taskkill", "/PID", fmt.Sprint(PID), "/F")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to close League client (PID: %d): %v", PID, err)
	}
	log.Printf("Closed League client (PID: %d)\n", PID)
	return nil
}

type Win32_Process struct {
	Name      string
	ProcessId uint32
}

// get the Process ID of the League Lobby window from its name
func getLobbyPID() (uint32, error) {
	var processes []Win32_Process
	query := fmt.Sprintf("SELECT Name, ProcessId FROM Win32_Process WHERE Name = '%s'", detection.LOBBY_WINDOW_NAME)

	// Retry loop with 5-second delay
	for {
		err := wmi.Query(query, &processes)
		if err == nil {
			if len(processes) > 0 {
				return processes[0].ProcessId, nil
			}
			return 0, fmt.Errorf("lobby window not found")
		}
		log.Printf("WMI query failed: %v", err)
		time.Sleep(5 * time.Second) // Delay before retrying
	}
}
