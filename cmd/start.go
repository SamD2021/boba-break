/*
Copyright © 2024 Samuel Dasilva
*/
package cmd

import (
	"time"

	"github.com/SamD2021/boba-break/tui/breakmanagerui"
	"github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// bm := new(breakmanager.BreakManager)
		workTime, err := cmd.Flags().GetString("work-duration")
		if err != nil {
			panic(err)
		}
		breakTime, _ := cmd.Flags().GetString("break-duration")
		if err != nil {
			panic(err)
		}
		// fmt.Printf("Work-duration: %s\nBreak-duration: %s\n", workTime, breakTime)
		workDuration, _ := time.ParseDuration(workTime)
		if err != nil {
			workDuration = time.Minute * 25
		}
		breakDuration, _ := time.ParseDuration(breakTime)
		if err != nil {
			breakDuration = time.Minute * 5
		}
		breakmanagerui.InitialModel(workDuration, breakDuration).Start()
	},
}

func init() {

	// Here you will define your flags and configuration settings.
	manageCmd.AddCommand(startCmd)
	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// startCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	startCmd.Flags().StringP("work-duration", "w", "25m", "Help message for duration")
	startCmd.Flags().StringP("break-duration", "b", "5m", "Help message for duration")
}
