package cmd

import (
	"time"

	"github.com/Link512/godid"
	"github.com/spf13/cobra"
)

var yesterdayCmd = &cobra.Command{
	Use:   "yesterday",
	Short: "Displays the tasks logged yesterday",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		yesterday, err := godid.GetYesterday()
		if err != nil {
			return err
		}
		if len(yesterday) != 0 {
			printResults(map[string][]string{time.Now().AddDate(0, 0, -1).Format("2006-01-02"): yesterday})
		} else {
			printEmpty()
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(yesterdayCmd)
}
