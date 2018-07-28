package cmd

import (
	"fmt"
	"os"
	"sort"
	"strings"

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
	bulkEntries := make([][]string, 0)
	for date, entries := range result {
		for _, entry := range entries {
			bulkEntries = append(bulkEntries, []string{date, entry})
		}
	}
	sort.Slice(bulkEntries, func(i, j int) bool {
		return strings.Compare(bulkEntries[i][0], bulkEntries[j][0]) > 0
	})
	writer.AppendBulk(bulkEntries)
	writer.Render()
}
