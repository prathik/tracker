package domain

import (
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

// PrintWeekData prints the data of the entire week in a tabular format
func PrintWeekData(sessionService *SessionService) ([][]string, error) {
	var weekData [][]string

	prevCount := 0.0
	prevTotal := 0.0
	duration, _ := time.ParseDuration("168h") // 7 days
	queryData, err := sessionService.ReportForPreviousDays(duration)

	sort.Slice(queryData, func(i, j int) bool {
		a, _ := time.Parse("2006-01-02", queryData[i].Day)
		b, _ := time.Parse("2006-01-02", queryData[j].Day)
		return a.After(b)
	})

	if err != nil {
		return nil, err
	}
	today := time.Now().Format("2006-01-02")
	for _, d := range queryData {
		day := d.Day
		if today != day {
			prevCount = prevCount + 1
			prevTotal = prevTotal + float64(d.Count)
		}
		weekData = append(weekData, []string{day, strconv.Itoa(d.Count)})
	}

	if prevCount != 0 {
		prevDaysAverage := int(prevTotal / prevCount)
		// TODO Convert to current week rolling average
		color.Green("Average Count = %d", prevDaysAverage)
	}
	
	return weekData, nil
}

// PrintByDay prints the data per day with items of each day in a table
func PrintByDay(ss *SessionService, since time.Duration) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Date", "Flow"})
	days, err := ss.ReportForPreviousDays(since)
	if err != nil {
		color.Red("error: %s", err)
		return
	}
	for _, d := range days {
		for _, wi := range d.Sessions {
			printData := []string{d.Day, strconv.Itoa(Score(wi.Challenge))}
			table.Rich(printData, []tablewriter.Colors{{}, getColour(Score(wi.Challenge))})
		}

	}
	table.SetAutoMergeCellsByColumnIndex([]int{0})
	table.Render()
}

func PrintWithTime(ss *SessionService, since time.Duration) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Time", "Delta", "Flow"})
	now := time.Now()
	days, err := ss.ReportForPreviousDays(since)
	if err != nil {
		color.Red("error: %s", err)
		return
	}
	for _, d := range days {
		for _, wi := range d.Sessions {
			printData := []string{wi.Time.Format(time.RFC3339), formatDelta(now, wi.Time), strconv.Itoa(Score(wi.Challenge))}
			table.Rich(printData, []tablewriter.Colors{{}, {}, getColour(Score(wi.Challenge))})
		}

	}
	table.Render()
}

func fmtDuration(d time.Duration) string {
	d = d.Round(time.Minute)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	return fmt.Sprintf("%dh%dm", h, m)
}

func formatDelta(now time.Time, doneOn time.Time) string {
	deltaInMinutes := now.Sub(doneOn).Minutes()
	if deltaInMinutes < 59 {
		return fmt.Sprintf("-%.0fm", now.Sub(doneOn).Minutes())
	}
	return "-" + fmtDuration(now.Sub(doneOn))
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

func HourAndMinuteFromString(str string) (hour int, minute int, err error) {
	if str == "" {
		err = errors.New("invalid input")
		return 0, 0, err
	}
	startTimeSplit := strings.SplitN(str, ":", 2)
	hour, err = strconv.Atoi(startTimeSplit[0])
	if err != nil {
		return 0, 0, err
	}
	minute, err = strconv.Atoi(startTimeSplit[1])
	if err != nil {
		return 0, 0, err
	}

	return hour, minute, err
}
