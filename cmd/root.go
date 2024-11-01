package cmd

import (
	"github.com/spf13/cobra"
	"notification-scheduler/executor"
	"notification-scheduler/types"
	"os"
)

var (
	mode  string
	input string
	title string
	body  string
	rrule string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "notification-scheduler",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		commandLineArguments := getCommandLineArguments()
		executor.Execute(commandLineArguments)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.notification-scheduler.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func getCommandLineArguments() types.CommandLineArguments {
	return types.CommandLineArguments{
		Mode:  mode,
		Input: input,
		Title: title,
		Body:  body,
		Rrule: rrule,
	}
}
