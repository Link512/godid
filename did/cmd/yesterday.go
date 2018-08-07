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
		godid.Init()
		defer godid.Close()
		yesterday, err := godid.GetYesterday()
		return handleResult(map[string][]string{time.Now().AddDate(0, 0, -1).Format("2006-01-02"): yesterday}, err)
	},
}

func init() {
	rootCmd.AddCommand(yesterdayCmd)
}
