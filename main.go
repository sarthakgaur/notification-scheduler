package main

import (
	"bufio"
	"database/sql"
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

func main() {
	database := getDatabase()
	defer database.Close()
	setupDatabase(database)

	rawNotificationSchedule := getInputFromStdin()
	cleanedNotificationSchedule := cleanNotificationSchedule(rawNotificationSchedule)
	validateNotificationSchedule(cleanedNotificationSchedule)
	createNotificationSchedule(database, cleanedNotificationSchedule)
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
		log.Fatalln("Error: Could not insert notification.")
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
