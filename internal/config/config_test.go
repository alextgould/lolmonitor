// Run using: go test ./internal/config

package config

import (
	"os"
	"testing"
	"time"
)

func TestSaveConfig(t *testing.T) {

	// Create a temporary file for testing
	tempFile, err := os.CreateTemp("", "config_test_*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	// Create a sample config
	cfg := Config{
		LoadOnStartup: true,
	}

	// Save the config
	if err := SaveConfig(tempFile.Name(), cfg); err != nil {
		t.Fatalf("SaveConfig failed: %v", err)
	}

	// Reload the config to verify it was saved correctly
	loadedCfg, err := LoadConfig(tempFile.Name())
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if loadedCfg.LoadOnStartup != cfg.LoadOnStartup {
		t.Errorf("Loaded config does not match saved config: got %+v, want %+v", loadedCfg, cfg)
	}
}

func TestLoadConfig(t *testing.T) {
	// Create a temporary test config file
	configContent := `{
		"dailyStartTime": "06:00",
		"gamesPerSession": 5
	}`

	tempFile := "test_config.json"
	err := os.WriteFile(tempFile, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create temp config file: %v", err)
	}
	defer os.Remove(tempFile) // Clean up after test

	cfg, err := LoadConfig(tempFile)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if cfg.GamesPerSession != 5 {
		t.Errorf("Expected 5 games per session, got %d", cfg.GamesPerSession)
	}

	if cfg.DailyStartTime != "06:00" {
		t.Errorf("Expected start time 04:00, got %s", cfg.DailyStartTime)
	}
}

func TestCheckConfig(t *testing.T) {
	// Create a temporary file for testing
	tempFile, err := os.CreateTemp("", "config_test_*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name()) // Clean up after test

	// Get the current time
	currentTime := time.Now()

	// Wait for a short duration to ensure a time difference
	time.Sleep(2 * time.Second)

	// Check if the file is detected as modified after a time in the past
	updated, err := CheckConfigUpdated(tempFile.Name(), currentTime.Add(-10*time.Second))
	if err != nil {
		t.Fatalf("CheckConfig failed: %v", err)
	}
	if !updated {
		t.Errorf("Expected file to be updated after the past time, but it was not")
	}

	// Check if the file is detected as not modified after a time in the future
	updated, err = CheckConfigUpdated(tempFile.Name(), currentTime.Add(10*time.Second))
	if err != nil {
		t.Fatalf("CheckConfig failed: %v", err)
	}
	if updated {
		t.Errorf("Expected file to not be updated after the future time, but it was")
	}

	// temp - force fail to manually review the log messages
	// t.Fail()
}

// manual test
// func TestLoadConfigManual(t *testing.T) {

// 	// load the config file in the build folder and display it

// 	cfg, err := LoadConfig("../../build/config.json") // this gives 30 OK
// 	if err != nil {
// 		t.Fatalf("LoadConfig failed: %v", err)
// 	}
// 	t.Logf("cfg.LobbyCloseDelaySeconds is %d", cfg.LobbyCloseDelaySeconds)

// 	// check if the config file has been updated in the last 60 seconds
// 	lastChecked := time.Now().Add(-60 * time.Second)
// 	configUpdated, err := CheckConfigUpdated("../../build/config.json", lastChecked)
// 	if err != nil {
// 		log.Printf("CheckConfigUpdated(\"\") error value is %v", err)
// 	}
// 	log.Printf("CheckConfigUpdated(\"\") result is %v", configUpdated)

// 	// works fine when using full path
// 	// Return value of fileInfo.ModTime().After(t) will be true, with fileInfo.ModTime() being 2025/04/18 15:09:15 and t value being 2025/04/18 15:08:20

// 	// now try using "" with the config file located in this directory

// 	cfg, err = LoadConfig("")
// 	if err != nil {
// 		t.Fatalf("LoadConfig failed: %v", err)
// 	}
// 	t.Logf("cfg.dailyStartTime is %v", cfg.DailyStartTime) // expect 2am

// 	// check if the config file has been updated in the last 60 seconds
// 	configUpdated, err = CheckConfigUpdated("", lastChecked)
// 	if err != nil {
// 		log.Printf("CheckConfigUpdated(\"\") error value is %v", err)
// 	}
// 	log.Printf("CheckConfigUpdated(\"\") result is %v", configUpdated)

// 	// works fine when using implied path
// 	// 2025/04/18 15:16:16 Return value of fileInfo.ModTime().After(t) will be true, with fileInfo.ModTime() being 2025/04/18 15:16:12 and t value being 2025/04/18 15:15:16

// 	t.Fail() // check log output
// }
