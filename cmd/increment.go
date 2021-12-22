package cmd

import (
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/prathik/tracker/domain"
	"github.com/prathik/tracker/repo"
	"github.com/spf13/cobra"
	"time"
)

// incrementCmd represents the increment command
var incrementCmd = &cobra.Command{
	Use:   "increment",
	Aliases: []string{"inc"},
	Short: "Increment number of sessions done",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		db := cmd.Flag("db").Value.String()
		bolt, err := repo.NewBoltDbRepo(db)
		if err != nil {
			color.Red("error: %s", err)
			return
		}
		defer bolt.Close()

		sessionService := domain.NewSessionService(bolt)

		challenge := args[0]

		var session *domain.Session
		session = &domain.Session{Challenge: challenge, Time: time.Now()}
		err = sessionService.Save(session)
		if err != nil {
			color.Red("error: %s", err)
			return
		}

		fmt.Printf("\n")
		PrintWeekData(sessionService)
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
}
