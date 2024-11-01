package cmd

import (
	"github.com/spf13/cobra"
	"notification-scheduler/executor"
)

var daemonCmd = &cobra.Command{
	Use:   "daemon",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		mode = "daemon"
		commandLineArguments := getCommandLineArguments()
		executor.Execute(commandLineArguments)
	},
}

func init() {
	rootCmd.AddCommand(daemonCmd)
}
