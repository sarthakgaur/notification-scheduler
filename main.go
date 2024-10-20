package main

import (
	"bufio"
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"strings"
	"time"
)

type notificationSchedule struct {
	title    string
	body     string
	datetime string
}

type commandLineArguments struct {
	mode     string
	input    string
	title    string
	body     string
	datetime string
}

func main() {
	cliArguments := getCommandLineArguments()
	database := getDatabase()
	defer database.Close()
	setupDatabase(database)
	executeProgramMode(cliArguments, database)
}

func getCommandLineArguments() commandLineArguments {
	mode := flag.String("mode", "schedule", `Mode of the program. Options: ["schedule", "daemon"]`)
	input := flag.String("input", "stdin", `Mode of input. Options: ["stdin", "cli"]`)
	title := flag.String("title", "", "Title of the notification")
	body := flag.String("body", "", "Body of the notification")
	datetime := flag.String("datetime", "", "Datetime of the notification in RFC3339 format")

	flag.Parse()

	return commandLineArguments{
		mode:     *mode,
		input:    *input,
		title:    *title,
		body:     *body,
		datetime: *datetime,
	}
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

func executeProgramMode(cliArguments commandLineArguments, database *sql.DB) {
	switch cliArguments.mode {
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

func executeScheduleMode(cliArguments commandLineArguments, database *sql.DB) {
	rawNotificationSchedule := getUserInput(cliArguments)
	cleanedNotificationSchedule := cleanNotificationSchedule(rawNotificationSchedule)
	validateNotificationSchedule(cleanedNotificationSchedule)
	createNotificationSchedule(database, cleanedNotificationSchedule)
}

func getUserInput(cliArguments commandLineArguments) notificationSchedule {
	switch cliArguments.input {
	case "stdin":
		return getInputFromStdin()
	case "cli":
		return getInputFromCommandLine(cliArguments)
	default:
		return getInputFromStdin()
	}
}

func getInputFromStdin() notificationSchedule {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter title: ")
	title, _ := reader.ReadString('\n')

	fmt.Print("Enter body: ")
	body, _ := reader.ReadString('\n')

	fmt.Print("Enter datetime (RFC3339): ")
	datetime, _ := reader.ReadString('\n')

	return notificationSchedule{
		title:    title,
		body:     body,
		datetime: datetime,
	}
}

func getInputFromCommandLine(cliArguments commandLineArguments) notificationSchedule {
	return notificationSchedule{
		title:    cliArguments.title,
		body:     cliArguments.body,
		datetime: cliArguments.datetime,
	}
}

func cleanNotificationSchedule(rawNotificationSchedule notificationSchedule) notificationSchedule {
	return notificationSchedule{
		title:    strings.TrimSpace(rawNotificationSchedule.title),
		body:     strings.TrimSpace(rawNotificationSchedule.body),
		datetime: strings.TrimSpace(rawNotificationSchedule.datetime),
	}
}

func validateNotificationSchedule(notification notificationSchedule) {
	if len(notification.title) == 0 {
		log.Fatalln("Error: Title is required.")
	}

	if len(notification.body) == 0 {
		log.Fatalln("Error: Body is required.")
	}

	if len(notification.datetime) == 0 {
		log.Fatalln("Error: Datetime is required.")
	}

	_, err := time.Parse(time.RFC3339, notification.datetime)
	if err != nil {
		log.Fatalln("Error: Datetime is not in RFC3339 format.", err)
	}
}

func createNotificationSchedule(database *sql.DB, notification notificationSchedule) {
	insertNotificationScheduleSQL := `INSERT INTO notification_schedule(created_on, title, body, datetime) VALUES (datetime('now'), ?, ?, ?)`

	statement, err := database.Prepare(insertNotificationScheduleSQL)
	if err != nil {
		log.Fatalln("Error: Could not prepare insert statement.")
	}

	title := notification.title
	body := notification.body
	datetime := notification.datetime

	_, err = statement.Exec(title, body, datetime)
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
		"datetime" TEXT NOT NULL
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
