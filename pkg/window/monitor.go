// monitor for windows opening or closing by name
package window

import (
	"log"
	"time"
)

// monitor for windows being opened or closed, sending ProcessEvent events to the events channel
// can pass 0 for delaySeconds and/or nil for pr to use default values (e.g. when testing)
// sample usage:
// myEvents := make(chan window.ProcessEvent, 10) // Buffered channel to prevent blocking
// go MonitorProcess(name="My window", myEvents, 0, nil) // Goroutine
// ...
// close(myEvents) // stop monitoring
func MonitorProcess(name string, events chan<- ProcessEvent, delaySeconds int, pr ProcessRetriever) {

	// defaults
	if pr == nil {
		pr = WMIProcessRetriever{}
	}
	if delaySeconds == 0 {
		delaySeconds = 1
	}

	wasActive := false
	for {
		isActive, id, _ := isProcessActive(name, pr)
		if !wasActive && isActive {
			events <- ProcessEvent{Name: name, PID: id, Type: "open", Timestamp: time.Now()}
		} else if wasActive && !isActive {
			events <- ProcessEvent{Name: name, PID: id, Type: "close", Timestamp: time.Now()}
		}
		time.Sleep(time.Duration(delaySeconds) * time.Second)
		wasActive = isActive
	}
}

// returns a bool for whether the process is active along with its process ID
// (this assumes there's only one instance and/or only the first instance is relevant)
// usage:
// isActive, PID, err = isProcessActive("my window", nil)
func isProcessActive(name string, pr ProcessRetriever) (bool, uint32, error) {

	processes, err := pr.GetProcessesByName(name)
	if err != nil {
		log.Printf("Error retrieving processes: %v", err)
		return false, 0, err
	}
	if len(processes) > 0 {
		return true, processes[0].ProcessId, nil
	}
	return false, 0, nil
}

// wait for a window to close
// can pass 0 for delaySeconds (default 1) and/or nil for pr (default WMI)
func WaitForProcessClose(name string, delaySeconds int, pr ProcessRetriever) {

	// defaults
	if pr == nil {
		pr = WMIProcessRetriever{}
	}
	if delaySeconds == 0 {
		delaySeconds = 1
	}

	for {
		isActive, _, _ := isProcessActive(name, pr)
		if !isActive {
			break
		}
		time.Sleep(time.Duration(delaySeconds) * time.Second)
	}
}
