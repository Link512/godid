package cmd

import (
	"fmt"
	"strings"

	"github.com/Link512/godid"

	"github.com/spf13/cobra"
)

var thisWeekCmd = &cobra.Command{
	Use:   "thisWeek",
	Short: "Display the tasks logged this week",
	RunE: func(cmd *cobra.Command, args []string) error {
		flat, err := cmd.Flags().GetBool("flat")
		if err != nil {
			return err
		}
		thisWeek, err := godid.GetThisWeek(flat)
		if err != nil {
			return err
		}
		for k, v := range thisWeek {
			fmt.Printf("Done on: %s\n%s\n", k, strings.Join(v, "\n"))
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(thisWeekCmd)
	thisWeekCmd.Flags().BoolP("flat", "f", false, "Do not aggregate the tasks per day")
}
