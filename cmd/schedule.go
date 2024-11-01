package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"notification-scheduler/executor"
)

var scheduleCmd = &cobra.Command{
	Use:   "schedule",
	Short: "Schedules a new notification task with specified parameters",
	Long:  getScheduleLongDescription(),
	Run: func(cmd *cobra.Command, args []string) {
		if input == "" {
			fmt.Println("Error: Please provide an input mode.")
			return
		}

		mode = "schedule"
		commandLineArguments := getCommandLineArguments()
		executor.Execute(commandLineArguments)
	},
}

func init() {
	rootCmd.AddCommand(scheduleCmd)

	scheduleCmd.Flags().StringVar(&input, "input", "stdin", "Input mode (stdin or cli)")
	scheduleCmd.Flags().StringVar(&title, "title", "", "Title of the schedule entry")
	scheduleCmd.Flags().StringVar(&body, "body", "", "Body/description of the schedule entry")
	scheduleCmd.Flags().StringVar(&rrule, "rrule", "", "Recurrence rule in iCalendar format (e.g., FREQ=WEEKLY)")
}

func getScheduleLongDescription() string {
	return `The "schedule" command allows you to create a new notification task that can be executed at a specified recurrence.
	
You can choose between two input modes:
- "stdin" mode to provide input directly from standard input.
- "cli" mode to provide input through command-line flags (--title, --body, and --rrule).

Examples:
1. Using CLI input mode with recurrence:
   notification-scheduler schedule --input cli --title "Meeting Reminder" --body "Don't forget the meeting at 10 AM" --rrule "FREQ=DAILY;INTERVAL=1"

2. Using stdin input mode:
   echo '{"title": "Task Reminder", "body": "Complete your tasks", "rrule": "FREQ=WEEKLY"}' | notification-scheduler schedule --input stdin`
}
