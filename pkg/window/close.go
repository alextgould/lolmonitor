// terminate windows by name
package window

import (
	"fmt"
)

// Close finds a process by name and closes it
func Close(name string, pr ProcessRetriever, pk ProcessKiller) error {
	isActive, PID, _ := isProcessActive(name, pr)
	if !isActive {
		return fmt.Errorf("process %v does not appear to be active", name)
	}
	return pk.KillProcess(PID)
}
