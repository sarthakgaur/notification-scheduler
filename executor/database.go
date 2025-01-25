package executor

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

func getNotificationSchedules(database *sql.DB) []notificationSchedule {
	getNotificationSchedulesSQL := `SELECT title, body, rrule FROM notification_schedule`

	rows, err := database.Query(getNotificationSchedulesSQL)
	if err != nil {
		log.Fatalln("Error: Could not query database.")
	}

	defer rows.Close()
	return transformNotificationScheduleRows(rows)
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

func transformNotificationScheduleRows(rows *sql.Rows) []notificationSchedule {
	var notificationSchedules []notificationSchedule

	for rows.Next() {
		var title string
		var body string
		var rrule string

		err := rows.Scan(&title, &body, &rrule)
		if err != nil {
			log.Fatalln("Error: Could not scan row.")
		}

		notificationSchedules = append(notificationSchedules, notificationSchedule{
			title: title,
			body:  body,
			rrule: rrule,
		})
	}

	return notificationSchedules
}
