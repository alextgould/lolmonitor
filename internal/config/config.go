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
	LobbyCloseDelaySeconds      int    `json:"lobbyCloseDelaySeconds"`
	LoadOnStartup               bool   `json:"loadOnStartup"`
	StartupInstalled            bool   `json:"startupInstalled"`
}

// DailyStartTime and DailyEndTime e.g. "04:00" "22:00"
var defaultConfig = Config{
	DailyStartTime:              "04:00",
	DailyEndTime:                "22:00",
	BreakBetweenGamesMinutes:    10,
	BreakBetweenSessionsMinutes: 60,
	GamesPerSession:             3,
	MinimumGameDurationMinutes:  15,
	LobbyCloseDelaySeconds:      10,
	LoadOnStartup:               true,
	StartupInstalled:            false,
}

func LoadConfig(filename string) (Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Config file not found, creating default config...")
		SaveConfig(filename, defaultConfig)
		return defaultConfig, nil
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	config := defaultConfig
	err = decoder.Decode(&config)
	if err != nil {
		fmt.Println("Invalid config file, resetting to default values...")
		SaveConfig(filename, defaultConfig)
		return defaultConfig, nil
	}

	return config, nil
}

// SaveConfig writes the configuration to the specified file.
func SaveConfig(filename string, cfg Config) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Pretty-print JSON
	return encoder.Encode(cfg)
}
