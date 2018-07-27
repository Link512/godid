package cmd

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
)

func printEmpty() {
	fmt.Println("Nothing here, you lazy slob!!")
}

func printResults(result map[string][]string) {
	if len(result) == 0 {
		printEmpty()
		return
	}
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
