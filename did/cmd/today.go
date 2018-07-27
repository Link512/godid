package cmd

import (
	"github.com/Link512/godid"

	"github.com/spf13/cobra"
)

var todayCmd = &cobra.Command{
	Use:   "today",
	Short: "Displays the tasks logged today",
	RunE: func(cmd *cobra.Command, args []string) error {
		today, err := godid.GetToday()
		if err != nil {
			return err
		}
		printResults(map[string][]string{"Today": today})
		return nil
	},
}

func init() {
	rootCmd.AddCommand(todayCmd)
}
