// load the config file, or create it if it's missing
package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/alextgould/lolmonitor/internal/interfaces/notifications"
	"github.com/alextgould/lolmonitor/internal/utils"
)

const CONFIG_FILE = "config.json"

type Config struct {
	BreakBetweenGamesMinutes    int    `json:"breakBetweenGamesMinutes"`
	BreakBetweenSessionsMinutes int    `json:"breakBetweenSessionsMinutes"`
	GamesPerSession             int    `json:"gamesPerSession"`
	MinimumGameDurationMinutes  int    `json:"minimumGameDurationMinutes"`
	LobbyCloseDelaySeconds      int    `json:"lobbyCloseDelaySeconds"`
	DailyStartTime              string `json:"dailyStartTime"`
	DailyEndTime                string `json:"dailyEndTime"`
	LoadOnStartup               bool   `json:"loadOnStartup"`
}

// DailyStartTime and DailyEndTime e.g. "04:00" "22:00"
var defaultConfig = Config{
	BreakBetweenGamesMinutes:    5,
	BreakBetweenSessionsMinutes: 60,
	GamesPerSession:             3,
	MinimumGameDurationMinutes:  15,
	LobbyCloseDelaySeconds:      30,
	DailyStartTime:              "00:00",
	DailyEndTime:                "00:00",
	LoadOnStartup:               true,
}

func SaveConfig(filename string, cfg Config) error {
	if filename == "" {
		filename = CONFIG_FILE
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Pretty-print JSON
	return encoder.Encode(cfg)
}

func LoadConfig(filename string) (Config, error) {
	if filename == "" {
		exePath, err := utils.GetCurrentPath()
		if err != nil {
			return Config{}, fmt.Errorf("failed to get executable path: %v", err)
		}
		filename = filepath.Join(exePath, CONFIG_FILE)
	}

	file, err := os.Open(filename)
	if err != nil {
		log.Println("Config file not found, creating default config...")
		SaveConfig(filename, defaultConfig)
		return defaultConfig, nil
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	config := defaultConfig       // create a config
	err = decoder.Decode(&config) // decode the JSON into the new config
	if err != nil {
		log.Printf("Error when decoding config file: %v", err)
		log.Println("Invalid config file, resetting to default values...")
		SaveConfig(filename, defaultConfig)
		return defaultConfig, nil
	}

	return config, nil
}

// Check if config file has been updated more recently than time value
func CheckConfigUpdated(filename string, t time.Time) (bool, error) {
	if filename == "" {
		filename = CONFIG_FILE
	}

	fileInfo, err := os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {

			// TEMP - getting default 30 sec delay when I have lobbyCloseDelaySeconds set to 15 sec in config - why??
			notification_text := fmt.Sprintf("File does not exist: %s", filename)
			notifications.SendNotification("Unable to find config file", notification_text, false)

			return false, fmt.Errorf("file does not exist: %s", filename)
		}
		return false, err
	}

	// Compare the file's modification time with the provided time
	return fileInfo.ModTime().After(t), nil
}
