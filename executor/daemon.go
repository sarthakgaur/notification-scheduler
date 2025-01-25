package executor

import (
	"github.com/gen2brain/beeep"
	_ "github.com/mattn/go-sqlite3"
	"github.com/teambition/rrule-go"
	"log"
	"time"
)

func executeDaemonMode() {
	for {
		dispatchPendingNotifications()
		nextMinute := time.Now().Truncate(time.Minute).Add(time.Minute)
		time.Sleep(time.Until(nextMinute))
	}
}

func dispatchPendingNotifications() {
	currentTime := time.Now().UTC().Truncate(time.Minute)
	log.Printf("Checking for pending notifications at %s.\n", currentTime)
	database := getDatabase()
	defer database.Close()
	notificationSchedules := getNotificationSchedules(database)
	matchingNotificationSchedules := getMatchingNotificationSchedules(notificationSchedules, currentTime)
	sendNotifications(matchingNotificationSchedules)
}

func getMatchingNotificationSchedules(notificationSchedules []notificationSchedule, currentTime time.Time) []notificationSchedule {
	var matchingNotificationSchedules []notificationSchedule

	for _, notificationSchedule := range notificationSchedules {
		rrule, err := rrule.StrToRRule(notificationSchedule.rrule)
		if err != nil {
			log.Fatalln("Error: Could not parse rrule.")
		}

		previousOccurrence := rrule.Before(currentTime, true)

		if isSameMinute(previousOccurrence, currentTime) {
			matchingNotificationSchedules = append(matchingNotificationSchedules, notificationSchedule)
		}
	}

	log.Printf("Found %d matching notification schedules.\n", len(matchingNotificationSchedules))
	return matchingNotificationSchedules
}

func sendNotifications(notificationSchedules []notificationSchedule) {
	for _, notificationSchedule := range notificationSchedules {
		err := beeep.Notify(notificationSchedule.title, notificationSchedule.body, "")
		if err != nil {
			log.Fatalln("Error: Could not send notification.")
		}
	}
}

func isSameMinute(time1 time.Time, time2 time.Time) bool {
	return time1.Year() == time2.Year() &&
		time1.Month() == time2.Month() &&
		time1.Day() == time2.Day() &&
		time1.Hour() == time2.Hour() &&
		time1.Minute() == time2.Minute()
}
