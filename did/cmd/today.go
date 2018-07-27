package cmd

import (
	"time"

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
		if len(today) != 0 {
			printResults(map[string][]string{time.Now().Format("2006-01-02"): today})
		} else {
			printEmpty()
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(todayCmd)
}
