package cmd

import (
	"github.com/Link512/godid"

	"github.com/spf13/cobra"
)

var thisWeekCmd = &cobra.Command{
	Use:   "thisWeek",
	Short: "Displays the tasks logged this week",
	RunE: func(cmd *cobra.Command, args []string) error {
		flat, err := cmd.Flags().GetBool("flat")
		if err != nil {
			return err
		}
		godid.Init()
		defer godid.Close()
		thisWeek, err := godid.GetThisWeek(flat)
		if err != nil {
			return err
		}
		printResults(thisWeek)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(thisWeekCmd)
	thisWeekCmd.Flags().BoolP("flat", "f", false, "Do not aggregate the tasks per day")
}
