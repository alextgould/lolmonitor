// internal/config/config.go
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type Config struct {
	DailyStartTime       string        `json:"dailyStartTime"`
	DailyEndTime         string        `json:"dailyEndTime"`
	BreakBetweenGames    time.Duration `json:"breakBetweenGames"`
	BreakBetweenSessions time.Duration `json:"breakBetweenSessions"`
	GamesPerSession      int           `json:"gamesPerSession"`
	MinGameDuration      time.Duration `json:"minGameDuration"`
}

var defaultConfig = Config{
	DailyStartTime:       "04:00",
	DailyEndTime:         "22:00",
	BreakBetweenGames:    15 * time.Minute,
	BreakBetweenSessions: 1 * time.Hour,
	GamesPerSession:      3,
	MinGameDuration:      15 * time.Minute,
}

func LoadConfig(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Config file not found, creating default config...")
		saveDefaultConfig(filename)
		return &defaultConfig, nil
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	config := defaultConfig
	err = decoder.Decode(&config)
	if err != nil {
		fmt.Println("Invalid config file, resetting to default values...")
		saveDefaultConfig(filename)
		return &defaultConfig, nil
	}

	return &config, nil
}

func saveDefaultConfig(filename string) {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("Failed to create default config file:", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(defaultConfig)
	if err != nil {
		fmt.Println("Failed to write default config file:", err)
	}
}
