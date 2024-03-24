package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/Link512/godid"
)

var rootCmd = &cobra.Command{
	Use:   "did",
	Short: "A simple task tracker",
	RunE: func(cmd *cobra.Command, args []string) error {
		entry, err := cmd.Flags().GetString("entry")
		if err != nil {
			return err
		}
		godid.Init()
		defer godid.Close()
		if entry != "" {
			return godid.AddEntry(entry)
		}
		reader := bufio.NewReader(os.Stdin)
		for {
			lineBytes, err := reader.ReadString('\n')
			line := strings.TrimSpace(string(lineBytes))
			if err != nil {
				if err == io.EOF && line != "" {
					return godid.AddEntry(line)
				}
				break
			}
			if err := godid.AddEntry(line); err != nil {
				return err
			}
		}
		return nil
	},
}

// Execute is the entry point for the CLI
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringP("entry", "e", "", "Entry to log")
}
