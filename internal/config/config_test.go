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
		"breakBetweenGamesMinutes": 15,
		"breakBetweenSessionsMinutes": 60,
		"gamesPerSession": 3,
		"minimumGameDurationMinutes": 15"
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

func TestSaveConfig(t *testing.T) {
	// Create a temporary file for testing
	tempFile, err := os.CreateTemp("", "config_test_*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	// Create a sample config
	cfg := Config{
		LoadOnStartup:    true,
		StartupInstalled: false,
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

	if loadedCfg.LoadOnStartup != cfg.LoadOnStartup || loadedCfg.StartupInstalled != cfg.StartupInstalled {
		t.Errorf("Loaded config does not match saved config: got %+v, want %+v", loadedCfg, cfg)
	}
}
