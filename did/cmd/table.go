package cmd

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/Link512/godid"

	"github.com/olekukonko/tablewriter"
)

func handleResult(result map[string][]string, err error) error {
	if err != nil {
		if didErr, ok := err.(godid.DidError); ok {
			return didErr
		}
		return errors.New("internal error, check the logs")
	}
	printResults(result)
	return nil
}

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
	writer.SetColWidth(4096)
	writer.SetHeader([]string{"Date", "Entries"})
	bulkEntries := make([][]string, 0)
	for date, entries := range result {
		for _, entry := range entries {
			bulkEntries = append(bulkEntries, []string{date, entry})
		}
	}
	if len(bulkEntries) == 0 {
		printEmpty()
		return
	}
	sort.Slice(bulkEntries, func(i, j int) bool {
		return strings.Compare(bulkEntries[i][0], bulkEntries[j][0]) < 0
	})
	writer.AppendBulk(bulkEntries)
	writer.Render()
}
