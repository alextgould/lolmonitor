// Contains the main types and wmi/exec functions used by the window package
// including both live (WMI) versions and Mock versions
package window

import (
	"fmt"
	"log"
	"os/exec"
	"time"

	"github.com/StackExchange/wmi"
)

// Win32_Process represents a Windows process
type Win32_Process struct {
	Name      string
	ProcessId uint32
}

// ProcessEvent captures open/close events
type ProcessEvent struct {
	Name      string
	PID       uint32
	Type      string // open, close
	Timestamp time.Time
}

// ProcessRetriever is an interface to allow mocking WMI process sourcing
type ProcessRetriever interface {
	GetProcessesByName(name string) ([]Win32_Process, error)
}

// ProcessRetriever will be WMIProcessRetriever when live
type WMIProcessRetriever struct{}

func (w WMIProcessRetriever) GetProcessesByName(name string) ([]Win32_Process, error) {
	var processes []Win32_Process
	query := fmt.Sprintf("SELECT Name, ProcessId FROM Win32_Process WHERE Name = '%s'", name)

	maxRetries := 5
	retryDelay := 5 * time.Second

	for i := 0; i < maxRetries; i++ {
		err := wmi.Query(query, &processes)
		if err == nil {
			return processes, nil
		}
		log.Printf("WMI query failed (attempt %d/%d): %v", i+1, maxRetries, err)
		time.Sleep(retryDelay)
	}

	return nil, fmt.Errorf("WMI query failed after %d retries", maxRetries)
}

// MockProcessRetriever is a mock implementation of ProcessRetriever for testing
type MockProcessRetriever struct {
	Processes map[string]Win32_Process
}

func (m *MockProcessRetriever) GetProcessesByName(name string) ([]Win32_Process, error) {
	var result []Win32_Process
	for _, process := range m.Processes {
		if process.Name == name {
			result = append(result, process)
		}
	}
	return result, nil
}

func (m *MockProcessRetriever) AddProcess(name string, pid uint32) {
	m.Processes[name] = Win32_Process{Name: name, ProcessId: pid}
}

func (m *MockProcessRetriever) RemoveProcess(name string) {
	delete(m.Processes, name)
}

// ProcessKiller is an interface to allow mocking task termination.

type ProcessKiller interface {
	KillProcess(pid uint32) error
}

// WMIProcessKiller implements real process termination.
type WMIProcessKiller struct{}

func (w WMIProcessKiller) KillProcess(pid uint32) error {
	cmd := exec.Command("taskkill", "/PID", fmt.Sprint(pid), "/F")
	return cmd.Run()
}

// MockProcessKiller implements process termination for testing.
type MockProcessKiller struct{}

func (m MockProcessKiller) KillProcess(pid uint32) error {
	return nil // Simulate successful termination
}
