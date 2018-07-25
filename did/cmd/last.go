package cmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Link512/godid"
	"github.com/spf13/cobra"
)

var lastCmd = &cobra.Command{
	Use:   "last",
	Short: "Display the entries logged in the last custom duration",
	Long:  `The duration string must be parsable by go's time.ParseDuration`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("must specify interval")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		flat, err := cmd.Flags().GetBool("flat")
		if err != nil {
			return err
		}
		last, err := godid.GetLastDuration(args[0], flat)
		if err != nil {
			return err
		}
		for k, v := range last {
			fmt.Printf("Done on: %s\n%s\n", k, strings.Join(v, "\n"))
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(lastCmd)
	lastCmd.Flags().BoolP("flat", "f", false, "Do not aggregate the tasks per day")
}
