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

type notificationSchedule struct {
	title string
	body  string
	rrule string
}

func Execute(cliArguments types.CommandLineArguments) {
	fmt.Println("Executing program with arguments:", cliArguments)
	database := getDatabase()
	defer database.Close()
	setupDatabase(database)
	executeProgramMode(cliArguments, database)
}

func getDatabase() *sql.DB {
	const NOTIFICATION_SCHEDULE_DB = "notification_schedule.db"

	db, err := sql.Open("sqlite3", NOTIFICATION_SCHEDULE_DB)
	if err != nil {
		log.Fatalln("Error: Could not open database.")
	}

	return db
}

func setupDatabase(database *sql.DB) {
	createNotificationScheduleTable(database)
	log.Println("Database ready.")
}

func executeProgramMode(cliArguments types.CommandLineArguments, database *sql.DB) {
	switch cliArguments.Mode {
	case "daemon":
		executeDaemonMode()
	case "schedule":
		executeScheduleMode(cliArguments, database)
	default:
		executeScheduleMode(cliArguments, database)
	}
}

func executeDaemonMode() {
	log.Println("Daemon mode not implemented.")
}

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

func createNotificationSchedule(database *sql.DB, notification notificationSchedule) {
	insertNotificationScheduleSQL := `INSERT INTO notification_schedule(created_on, title, body, rrule) VALUES (datetime('now'), ?, ?, ?)`

	statement, err := database.Prepare(insertNotificationScheduleSQL)
	if err != nil {
		log.Fatalln("Error: Could not prepare insert statement.")
	}

	title := notification.title
	body := notification.body
	rrule := notification.rrule

	_, err = statement.Exec(title, body, rrule)
	if err != nil {
		log.Fatalln("Error: Could not insert notification schedule.")
	}

	log.Println("Notification schedule created.")
}

func createNotificationScheduleTable(database *sql.DB) {
	createNotificationScheduleTableSQL := `CREATE TABLE IF NOT EXISTS notification_schedule (
		"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,		
		"created_on" TEXT NOT NULL,
		"title" TEXT NOT NULL,
		"body" TEXT NOT NULL,
		"rrule" TEXT NOT NULL
	  );`

	statement, err := database.Prepare(createNotificationScheduleTableSQL)
	if err != nil {
		log.Fatalln("Error: Could not prepare create table statement.")
	}

	_, err = statement.Exec()
	if err != nil {
		log.Fatalln("Error: Could not create table.")
	}
}
