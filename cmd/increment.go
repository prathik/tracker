package cmd

import (
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

func init() {
	rootCmd.AddCommand(incrementCmd)
}
