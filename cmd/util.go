package cmd

import (
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/prathik/tracker/service"
	"os"
	"strconv"
)

// PrintWeekData prints the data of the entire week in a tabular format
func PrintWeekData(ss *service.SessionService) {
	// Print entire week as a table for a reminder
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Day", "Count"})

	for _, d := range ss.QueryData(7).DayDataCollection {
		table.Append([]string{d.Time.Format("2006-01-02"), strconv.Itoa(d.Count)})
	}
	table.Render()
}

// PrintByDay prints the data per day with items of each day in a table
func PrintByDay(ss *service.SessionService) {
	for _, d := range ss.QueryData(7).DayDataCollection {
		color.Green("%s", d.Time.Format("2006-01-02"))
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Work", "Joy", "Importance", "Notes"})
		for _, wi := range d.WorkItem {
			printData := []string{wi.Work, strconv.Itoa(wi.Joy), strconv.Itoa(wi.Importance), wi.Notes}
			table.Rich(printData, []tablewriter.Colors{{}, getColour(wi.Joy), getColour(wi.Importance), {}})
		}
		table.Render()
	}
}

func getColour(value int) tablewriter.Colors {
	if value > 6 {
		return tablewriter.Colors{tablewriter.Bold, tablewriter.FgCyanColor}
	}

	if value < 4 {
		return tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiRedColor}
	}

	return tablewriter.Colors{tablewriter.Normal, tablewriter.ALIGN_DEFAULT}
}
