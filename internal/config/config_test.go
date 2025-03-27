// Run using: go test ./internal/config

package config

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// Create a temporary test config file
	configContent := `{
		"dailyStartTime": "04:00",
		"dailyEndTime": "22:00",
		"breakBetweenGames": "00:15",
		"breakBetweenSessions": "01:00",
		"gamesPerSession": 3,
		"minGameDuration": "00:15"
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

	if cfg.GamesPerSession != 3 {
		t.Errorf("Expected 3 games per session, got %d", cfg.GamesPerSession)
	}

	if cfg.DailyStartTime != "04:00" {
		t.Errorf("Expected start time 04:00, got %s", cfg.DailyStartTime)
	}
}
