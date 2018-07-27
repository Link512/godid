package cmd

import (
	"os"

	"github.com/olekukonko/tablewriter"
)

func printResults(result map[string][]string) {
	writer := tablewriter.NewWriter(os.Stdout)
	writer.SetAutoMergeCells(true)
	writer.SetAutoWrapText(true)
	writer.SetRowLine(true)
	writer.SetHeader([]string{"Date", "Entries"})
	for date, entries := range result {
		for _, entry := range entries {
			writer.Append([]string{date, entry})
		}
	}
	writer.Render()
}
