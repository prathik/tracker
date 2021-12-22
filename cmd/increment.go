package cmd

import (
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/prathik/tracker/repo"
	"github.com/prathik/tracker/service"
	"github.com/spf13/cobra"
	"strconv"
	"time"
)

// incrementCmd represents the increment command
var incrementCmd = &cobra.Command{
	Use:   "increment",
	Short: "Increment number of sessions done",
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		db := cmd.Flag("db").Value.String()
		bolt := repo.NewBoltDbRepo(db)
		defer bolt.Close()
		ss := service.NewSessionService(bolt)


		joy, _ := strconv.Atoi(args[0])
		impact, _ := strconv.Atoi(args[1])

		startTime, _ := cmd.Flags().GetString("start-time")
		count, _ := cmd.Flags().GetInt("count")

		const (
			notesResult = "deprecated"
		)

		// Create count number of sessions
		if startTime == "" {
			if count > 1 {
				color.Red("count > 1 is only supported when --start-time flag is passed")
				return
			}
			ss.Create(&service.Item{Joy: joy, Impact: impact, Notes: notesResult, Time: time.Now()})
		} else {
			for i := 0; i < count; i++ {
				sessionTime, err := SessionTime(startTime, i)
				if err != nil {
					color.Red(err.Error())
					return
				}
				ss.Create(&service.Item{Joy: joy, Impact: impact, Notes: notesResult, Time: sessionTime})
			}
		}

		fmt.Printf("\n")

		PrintWeekData(ss)
	},
}

func SessionTime(startTime string, count int) (time.Time, error) {
	hour, min, err := HourAndMinuteFromString(startTime)
	if err != nil {
		return time.Time{}, err
	}
	now := time.Now()
	sessionTime := time.Date(now.Year(), now.Month(), now.Day(), hour, min, 0, 0, now.Location())
	totalMinutes := time.Duration(30*count) * time.Minute
	sessionTime = sessionTime.Add(totalMinutes)

	if now.Before(sessionTime) {
		return time.Time{}, errors.New("session time ahead of current time")
	}
	return sessionTime, nil
}

func init() {
	rootCmd.AddCommand(incrementCmd)
	incrementCmd.Flags().String("start-time", "", "Start hour and minute, use hh:mm format")
	incrementCmd.Flags().Int("count", 1, "Count of sessions done")
}
