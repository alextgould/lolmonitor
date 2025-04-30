package application

import (
	"log"
	"testing"
	"time"

	"github.com/alextgould/lolmonitor/internal/config"
)

func TestPostGame(t *testing.T) {
	cfg := config.Config{
		MinimumGameDurationMinutes:  10,
		GamesPerSession:             3,
		BreakBetweenSessionsMinutes: 30,
		BreakBetweenGamesMinutes:    5,
	}

	var gameStartTime, gameEndTime, currentTime, endOfBreak time.Time
	var breakDuration time.Duration
	var sessionGames, newSessionGames int

	sessionGames = 0
	currentTime = time.Now()
	gameStartTime = time.Now().Add(-5 * time.Minute)
	gameEndTime = time.Now()

	// check remakes don't count
	newSessionGames, endOfBreak, _ = postGame(cfg, time.Now(), gameStartTime, gameEndTime, sessionGames)
	if newSessionGames != 0 {
		t.Errorf("Expected sessionGames to be 0, got %d", newSessionGames)
	}
	if endOfBreak.After(currentTime) {
		t.Errorf("Expected endOfBreak to be currentTime, got %v and %v respectively", endOfBreak, currentTime)
	}

	// check regular games do count
	gameStartTime = time.Now().Add(-20 * time.Minute)
	newSessionGames, _, breakDuration = postGame(cfg, time.Now(), gameStartTime, gameEndTime, sessionGames)
	if newSessionGames != 1 {
		t.Errorf("Expected sessionGames to be 1, got %d", newSessionGames)
	}
	if breakDuration != 5*time.Minute {
		t.Errorf("Expected short breakDuration of 5 minutes, got %v", breakDuration)
	}

	// check end of session logic
	sessionGames = 2
	newSessionGames, _, breakDuration = postGame(cfg, time.Now(), gameStartTime, gameEndTime, sessionGames)
	if newSessionGames != 0 {
		t.Errorf("Expected sessionGames to be 0, got %d", newSessionGames)
	}
	if breakDuration != 30*time.Minute {
		t.Errorf("Expected breakDuration to be 30 minutes, got %v", breakDuration)
	}
}

func TestIsLobbyBan(t *testing.T) {
	cfg := config.Config{
		DailyStartTime: "04:00",
		DailyEndTime:   "22:00",
	}

	var currentTime, endOfBreak time.Time
	log.Println("New test")
	// test lobby can open during a normal time
	currentTime = time.Date(2023, 10, 1, 10, 0, 0, 0, time.UTC) // 10:00 AM
	endOfBreak = currentTime

	isBanned, err := isLobbyBan(cfg, currentTime, endOfBreak)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if isBanned {
		t.Errorf("Expected lobby to not be banned")
	}
	log.Println("New test")

	// test before daily start time
	currentTime = time.Date(2023, 10, 1, 3, 0, 0, 0, time.UTC) // 3:00 AM
	endOfBreak = currentTime

	isBanned, err = isLobbyBan(cfg, currentTime, endOfBreak)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if !isBanned {
		t.Errorf("Expected lobby to be banned due to DailyStartTime of 3am, but it was not")
	}
	log.Println("New test")
	// test after daily end time
	currentTime = time.Date(2023, 10, 1, 23, 0, 0, 0, time.UTC) // 11:00 PM
	endOfBreak = currentTime
	isBanned, err = isLobbyBan(cfg, currentTime, endOfBreak)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if !isBanned {
		t.Errorf("Expected lobby to be banned due to DailyEndTime of 11pm, but it was not")
	}
	log.Println("New test")
	// test lobby can open if DailyStartTime is "00:00"
	cfg = config.Config{
		DailyStartTime: "00:00",
		DailyEndTime:   "00:00",
	}
	isBanned, err = isLobbyBan(cfg, currentTime, endOfBreak)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if isBanned {
		t.Errorf("Expected lobby to not be banned due to 00:00 DailyStartTime, but it was banned")
	}
	log.Println("New test")
	// test lobby can open if DailyEndTime is "00:00"
	currentTime = time.Date(2023, 10, 1, 23, 0, 0, 0, time.UTC) // 11:00 PM
	isBanned, err = isLobbyBan(cfg, currentTime, endOfBreak)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if isBanned {
		t.Errorf("Expected lobby to not be banned due to 00:00 DailyEndTime, but it was banned")
	}
	log.Println("New test")
	// test lobby can't open if during a break
	currentTime = time.Date(2023, 10, 1, 10, 0, 0, 0, time.UTC) // 10:00 AM
	endOfBreak = currentTime.Add(10 * time.Minute)

	isBanned, err = isLobbyBan(cfg, currentTime, endOfBreak)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if !isBanned {
		t.Errorf("Expected lobby to be banned due to endOfBreak")
	}
}

// manual test to replicate notification after hours
// actually this all seems to be working fine
// func TestPostGameAfterHours(t *testing.T) {
// 	cfg := config.Config{
// 		MinimumGameDurationMinutes:  10,
// 		GamesPerSession:             3,
// 		BreakBetweenSessionsMinutes: 30,
// 		BreakBetweenGamesMinutes:    120,
// 		DailyEndTime:                "01:00",
// 	}

// 	sessionGames := 0
// 	// currentTime := time.Date(2023, 10, 1, 23, 0, 0, 0, time.UTC) // 11:00 PM
// 	currentTime := time.Now()
// 	gameStartTime := currentTime.Add(-45 * time.Minute)
// 	gameEndTime := currentTime

// 	newSessionGames, endOfBreak, breakDuration := postGame(cfg, currentTime, gameStartTime, gameEndTime, sessionGames)
// 	log.Printf("newSessionGames is %v", newSessionGames)
// 	log.Printf("endOfBreak is %v", endOfBreak)
// 	log.Printf("breakDuration is %v", breakDuration)
// 	// display the lobby blocked notification using this endOfBreak value
// 	// says "try again in -830219 minutes at 11:06pm"
// 	// suggests an issue with math.Ceil(time.Until(endOfBreak).Minutes())
// 	notifications.LobbyBlocked(endOfBreak)
// 	t.Fail() // show time values in log
// 	// 2025/04/30 22:04:49 endOfBreak is 2023-10-01 23:05:00 +0000 UTC
// 	// 2025/04/30 22:04:49 breakDuration is 5m0s
// }
