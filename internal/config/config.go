// load the config file, or create it if it's missing
package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/alextgould/lolmonitor/internal/utils"
)

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

// if filename is not explicitly specified (i.e. ""), use the executable path and "config.json"
func defaultPath(filename string) (string, error) {
	if filename == "" {
		exePath, err := utils.GetCurrentPath(false)
		if err != nil {
			return "", err
		}
		filename = filepath.Join(exePath, "config.json")
	}
	return filename, nil
}

func SaveConfig(filename string, cfg Config) error {

	// use default path and filename unless specified explicitly
	filename, err := defaultPath(filename)
	if err != nil {
		return err
	}

	// create config file
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// save cfg in json format
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Pretty-print JSON
	return encoder.Encode(cfg)
}

func LoadConfig(filename string) (Config, error) {

	// use default path and filename unless specified explicitly
	filename, err := defaultPath(filename)
	if err != nil {
		return Config{}, err
	}

	// open file, creating a new config file if one is not found
	file, err := os.Open(filename)
	if err != nil {
		log.Println("Config file not found, creating default config...")
		SaveConfig(filename, defaultConfig)
		return defaultConfig, nil
	}
	defer file.Close()

	// load json into config
	decoder := json.NewDecoder(file)
	config := defaultConfig
	err = decoder.Decode(&config)
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

	// use default path and filename unless specified explicitly
	filename, err := defaultPath(filename)
	if err != nil {
		return false, err
	}

	// check if date modified is after time t
	fileInfo, err := os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return false, fmt.Errorf("unable to find config file: %s", filename)
		}
		return false, err
	}

	// Compare the file's modification time with the provided time
	return fileInfo.ModTime().After(t), nil
}
