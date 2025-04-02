// load the config file, or create it if it's missing
package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	DailyStartTime              string `json:"dailyStartTime"`
	DailyEndTime                string `json:"dailyEndTime"`
	BreakBetweenGamesMinutes    int    `json:"breakBetweenGamesMinutes"`
	BreakBetweenSessionsMinutes int    `json:"breakBetweenSessionsMinutes"`
	GamesPerSession             int    `json:"gamesPerSession"`
	MinimumGameDurationMinutes  int    `json:"minimumGameDurationMinutes"`
}

var defaultConfig = Config{
	DailyStartTime:              "04:00",
	DailyEndTime:                "22:00",
	BreakBetweenGamesMinutes:    15,
	BreakBetweenSessionsMinutes: 60,
	GamesPerSession:             3,
	MinimumGameDurationMinutes:  15,
}

func LoadConfig(filename string) (Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Config file not found, creating default config...")
		createDefaultConfig(filename)
		return defaultConfig, nil
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	config := defaultConfig
	err = decoder.Decode(&config)
	if err != nil {
		fmt.Println("Invalid config file, resetting to default values...")
		createDefaultConfig(filename)
		return defaultConfig, nil
	}

	return config, nil
}

func createDefaultConfig(filename string) {
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
