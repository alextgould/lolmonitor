// send windows toast notifications
// Note these will not appear if you have notifications turned off in Settings > System > Notifications & actions > Get notifications from apps and other senders

package notifications

import (
	"fmt"
	"log"
	"math"
	"time"

	"github.com/go-toast/toast"
)

func SendNotification(title, message string, duration_long bool) { // TODO: add optional button text
	// short duration is ~7 sec, long duration is ~25 sec
	duration := toast.Short
	if duration_long {
		duration = toast.Long
	}
	notification := toast.Notification{
		AppID:   "LoL Monitor",
		Title:   title,
		Message: message,
		// Icon: "go.png",  // TODO: create an icon for the program and add it here
		// Actions: []toast.Action{  // This file must exist (remove this line if it doesn't)  // TODO: have optional actions, do I need separate functions?
		//     {"protocol", "I'm a button", ""},
		// },
		Duration: duration,
	}
	err := notification.Push()
	if err != nil {
		log.Printf("Failed to send notification: %v", err)
	}
}

func DelayClose(delay int, sessionGames, gamesPerSession int) {
	var notification_title, notification_text string
	if sessionGames > 0 {
		notification_title = fmt.Sprintf("Game over! That's game %d of %d.", sessionGames, gamesPerSession)
	} else { // sessionGames is reset to 0 after the last game of the session
		notification_title = fmt.Sprintf("Session over! That's game %d of %d.", gamesPerSession, gamesPerSession)
	}
	notification_text = fmt.Sprintf("The lobby will close automatically in %d seconds", delay)
	SendNotification(notification_title, notification_text, false)
}

func EndOfGame(endOfBreak time.Time, sessionGames, gamesPerSession int) {
	var notification_title, notification_text string
	if sessionGames > 0 {
		notification_title = fmt.Sprintf("Game over! That's game %d of %d.", sessionGames, gamesPerSession)
		notification_text = fmt.Sprintf("Take a short break to reset and come back in %.0f minutes at %v",
			math.Ceil(time.Until(endOfBreak).Minutes()),
			endOfBreak.Format("3:04pm"))
	} else {
		notification_title = fmt.Sprintf("Session over! That's game %d of %d.", gamesPerSession, gamesPerSession)
		notification_text = fmt.Sprintf("Take a break to recover and come back some time after %v)",
			endOfBreak.Format("3:04pm"))
	}
	SendNotification(notification_title, notification_text, true)
}

func LobbyBlocked(endOfBreak time.Time) {
	notification_title := "You're on a break!"
	notification_text := fmt.Sprintf("Try again in %.0f minutes at %v",
		math.Ceil(time.Until(endOfBreak).Minutes()),
		endOfBreak.Format("3:04pm"))
	SendNotification(notification_title, notification_text, false)
}
