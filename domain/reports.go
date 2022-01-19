package domain

import (
	"strconv"
	"time"
)

type weekData struct {
	Raw          [][]string
	AverageCount float64
	FlowRatio    float64
}

// GenerateWeekReport generates the data of the entire week in a tabular format
func GenerateWeekReport(now time.Time, sessionService *SessionService) (*weekData, error) {
	var data weekData

	duration := timeToStartOfTheWeek(now)
	queryData, err := sessionService.ReportForPreviousDays(duration)

	SortDays(queryData)

	dayCount := 0.0
	total := 0.0

	if err != nil {
		return nil, err
	}

	flowCount := 0.0

	for _, d := range queryData {
		day := d.Day
		dayCount = dayCount + 1
		for _, session := range d.Sessions {
			if session.Challenge == "flow" {
				flowCount++
			}
		}
		total = total + float64(d.Count)
		data.Raw = append(data.Raw, []string{day, strconv.Itoa(d.Count)})
	}

	averagePomodoroPerDay := total / dayCount
	data.AverageCount = averagePomodoroPerDay
	data.FlowRatio = flowCount / total

	return &data, nil
}

func timeToStartOfTheWeek(current time.Time) time.Duration {
	days := (mod(int(current.Weekday()-1), 7) * 24) + current.Hour()
	duration, _ := time.ParseDuration(strconv.Itoa(days) + "h") // 7 days
	return duration
}

