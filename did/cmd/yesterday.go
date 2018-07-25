package cmd

import (
	"fmt"
	"strings"

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
		fmt.Println(strings.Join(yesterday, "\n"))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(yesterdayCmd)
}
