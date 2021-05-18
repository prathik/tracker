package cmd

import (
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/prathik/tracker/service"
	"os"
	"strconv"
	"time"
)

// PrintWeekData prints the data of the entire week in a tabular format
func PrintWeekData(ss *service.SessionService) {
	// Print entire week as a table for a reminder
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Day", "Count"})
	prevCount := 0.0
	prevTotal := 0.0
	queryData := ss.QueryData(7)
	today := time.Now().Format("2006-01-02")
	for _, d := range queryData.DayDataCollection {
		day := d.Time.Format("2006-01-02")
		if today != day {
			prevCount = prevCount + 1
			prevTotal = prevTotal + float64(d.Count)
		}
		table.Append([]string{day, strconv.Itoa(d.Count)})
	}
	prevDaysAverage := int(prevTotal / prevCount)
	color.Green("Count to meet or exceed today = %d", prevDaysAverage)
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
