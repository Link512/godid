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
		godid.Init()
		defer godid.Close()
		today, err := godid.GetToday()
		return handleResult(map[string][]string{time.Now().Format("2006-01-02"): today}, err)
	},
}

func init() {
	rootCmd.AddCommand(todayCmd)
}
