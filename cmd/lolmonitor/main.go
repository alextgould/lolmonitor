// To run (in terminal, at project root): go run cmd/lolmonitor/main.go
// To build: go build -o build/lolmonitor.exe -ldflags "-H windowsgui" cmd/lolmonitor/main.go
// -ldflags "-H windowsgui"   -- hides the terminal
// To build a version that will run with output logged to the terminal: go build -o build/lolmonitor_terminal.exe cmd/lolmonitor/main.go

package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/alextgould/lolmonitor/internal/application"
	"github.com/alextgould/lolmonitor/internal/config"
	"github.com/alextgould/lolmonitor/internal/infrastructure/startup"

	window "github.com/alextgould/lolmonitor/pkg/window"
)

func main() {
	log.Println("Starting up lolmonitor")

	// Load (or create) config file
	cfg, err := config.LoadConfig("")
	if err != nil {
		log.Panicf("Failed to load config: %v", err)
	}

	// Add or remove startup process based on config settings
	if cfg.LoadOnStartup {
		err = startup.ConfirmLoadOnStartup()
	} else {
		err = startup.ConfirmNoLoadOnStartup()
	}
	if err != nil {
		log.Panicf("Error updating startup: %v", err)
	}

	// Create channel for window events
	gameEvents := make(chan window.ProcessEvent, 10) // Buffered channel to prevent blocking

	// Start monitoring in a separate goroutine
	go application.Monitor(cfg, gameEvents, nil, nil)

	// Handle Ctrl+C for graceful shutdown
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	<-signalChan // Block until an interrupt signal is received (or else indefinitely)
	log.Println("Shutting down lolmonitor")
	close(gameEvents) // Signal monitoring to stop
}
