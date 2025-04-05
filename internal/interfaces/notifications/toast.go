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

func SendNotification(title, message string) { // TODO: add optional button text
	notification := toast.Notification{
		AppID:   "LoL Monitor",
		Title:   title,
		Message: message,
		// Icon: "go.png",  // TODO: create an icon for the program and add it here
		// Actions: []toast.Action{  // This file must exist (remove this line if it doesn't)  // TODO: have optional actions, do I need separate functions?
		//     {"protocol", "I'm a button", ""},
		// },
	}
	err := notification.Push()
	if err != nil {
		log.Printf("Failed to send notification: %v", err)
	}
}

func EndOfGame(endOfBreak time.Time) {
	notification_text := fmt.Sprintf("GG. Take a break and come back in %.0f minutes at %v",
		math.Ceil(time.Until(endOfBreak).Minutes()),
		endOfBreak.Format("3:04pm"))
	SendNotification("League of Legends Lobby Closed", notification_text)
}

func LobbyBlocked(endOfBreak time.Time) {
	notification_text := fmt.Sprintf("You're on a break! Try again in %.0f minutes at %v",
		math.Ceil(time.Until(endOfBreak).Minutes()),
		endOfBreak.Format("3:04pm"))
	SendNotification("League of Legends Blocked", notification_text)
}
