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

type notification struct {
	title    string
	body     string
	datetime string
}

func main() {
	database := getDatabase()
	defer database.Close()
	setupDatabase(database)

	rawNotificationInput := getInputFromStdin()
	notificationInput := cleanNotificationInput(rawNotificationInput)
	checkNotificationInputs(notificationInput)
	insertNotification(database, notificationInput)
}

func getDatabase() *sql.DB {
	const fileName = "sqlite.db"

	db, err := sql.Open("sqlite3", fileName)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func setupDatabase(database *sql.DB) {
	createNotificationTable(database)
	log.Println("Database ready.")
}

func getInputFromStdin() notification {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter title: ")
	title, _ := reader.ReadString('\n')

	fmt.Print("Enter body: ")
	body, _ := reader.ReadString('\n')

	fmt.Print("Enter datetime (RFC3339): ")
	datetime, _ := reader.ReadString('\n')

	return notification{
		title:    title,
		body:     body,
		datetime: datetime,
	}
}

func cleanNotificationInput(rawNotification notification) notification {
	return notification{
		title:    strings.TrimSpace(rawNotification.title),
		body:     strings.TrimSpace(rawNotification.body),
		datetime: strings.TrimSpace(rawNotification.datetime),
	}
}

func checkNotificationInputs(notification notification) {
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

func insertNotification(database *sql.DB, notification notification) {
	insertNotificationSQL := `INSERT INTO notification(created_on, title, body, datetime) VALUES (datetime('now'), ?, ?, ?)`
	statement, err := database.Prepare(insertNotificationSQL)
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
}

func createNotificationTable(database *sql.DB) {
	createNotificationTableSQL := `CREATE TABLE IF NOT EXISTS notification (
		"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,		
		"created_on" TEXT NOT NULL,
		"title" TEXT NOT NULL,
		"body" TEXT NOT NULL,
		"datetime" TEXT NOT NULL
	  );`

	statement, err := database.Prepare(createNotificationTableSQL)
	if err != nil {
		log.Fatalln("Error: Could not prepare create table statement.")
	}

	_, err = statement.Exec()
	if err != nil {
		log.Fatalln("Error: Could not create table.")
	}
}
