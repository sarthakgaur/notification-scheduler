package executor

import (
	"bufio"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/teambition/rrule-go"
	"log"
	"notification-scheduler/types"
	"os"
	"strings"
)

func executeScheduleMode(cliArguments types.CommandLineArguments, database *sql.DB) {
	rawNotificationSchedule := getUserInput(cliArguments)
	cleanedNotificationSchedule := cleanNotificationSchedule(rawNotificationSchedule)
	validateNotificationSchedule(cleanedNotificationSchedule)
	createNotificationSchedule(database, cleanedNotificationSchedule)
}

func getUserInput(cliArguments types.CommandLineArguments) notificationSchedule {
	switch cliArguments.Input {
	case "stdin":
		return getInputFromStdin()
	case "cli":
		return getInputFromCommandLine(cliArguments)
	default:
		fmt.Println("default case")
		return getInputFromStdin()
	}
}

func getInputFromStdin() notificationSchedule {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter title: ")
	title, _ := reader.ReadString('\n')

	fmt.Print("Enter body: ")
	body, _ := reader.ReadString('\n')

	fmt.Print("Enter recurrence rule (RFC5545): ")
	rrule, _ := reader.ReadString('\n')

	return notificationSchedule{
		title: title,
		body:  body,
		rrule: rrule,
	}
}

func getInputFromCommandLine(cliArguments types.CommandLineArguments) notificationSchedule {
	return notificationSchedule{
		title: cliArguments.Title,
		body:  cliArguments.Body,
		rrule: cliArguments.Rrule,
	}
}

func cleanNotificationSchedule(rawNotificationSchedule notificationSchedule) notificationSchedule {
	return notificationSchedule{
		title: strings.TrimSpace(rawNotificationSchedule.title),
		body:  strings.TrimSpace(rawNotificationSchedule.body),
		rrule: strings.TrimSpace(rawNotificationSchedule.rrule),
	}
}

func validateNotificationSchedule(notification notificationSchedule) {
	if len(notification.title) == 0 {
		log.Fatalln("Error: Title is required.")
	}

	if len(notification.body) == 0 {
		log.Fatalln("Error: Body is required.")
	}

	if len(notification.rrule) == 0 {
		log.Fatalln("Error: Recurrence rule is required.")
	}

	_, err := rrule.StrToRRule(notification.rrule)
	if err != nil {
		log.Fatalln("Error: Invalid recurrence rule.")
	}
}
