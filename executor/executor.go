package executor

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"notification-scheduler/types"
)

type notificationSchedule struct {
	title string
	body  string
	rrule string
}

func Execute(cliArguments types.CommandLineArguments) {
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
