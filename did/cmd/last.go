package cmd

import (
	"errors"

	"github.com/Link512/godid"
	"github.com/spf13/cobra"
)

var lastCmd = &cobra.Command{
	Use:   "last",
	Short: "Display the tasks logged in the last custom day duration",
	Long:  `The duration string only accepts days: 1d, 2d, etc`,
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
		godid.Init()
		defer godid.Close()
		last, err := godid.GetLastDuration(args[0], flat)
		return handleResult(last, err)
	},
}

func init() {
	rootCmd.AddCommand(lastCmd)
	lastCmd.Flags().BoolP("flat", "f", false, "Do not aggregate the tasks per day")
}
