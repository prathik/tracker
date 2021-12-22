package cmd

import (
	"github.com/fatih/color"
	"github.com/prathik/tracker/domain"
	"github.com/prathik/tracker/repo"
	"github.com/spf13/cobra"
	"time"
)

var daysSince int

// showCmd represents the show command
var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Shows the activities done in the past (by default week)",
	Run: func(cmd *cobra.Command, args []string) {
		db := cmd.Flag("db").Value.String()
		bolt, err := repo.NewBoltDbRepo(db)
		if err != nil {
			color.Red("error: %s", err)
			return
		}
		defer bolt.Close()
		ss := domain.NewSessionService(bolt)

		includeTime, _ := cmd.Flags().GetBool("with-time")

		if includeTime {
			PrintWithTime(ss, time.Duration(daysSince*24)*time.Hour)
		} else {
			PrintByDay(ss, time.Duration(daysSince*24)*time.Hour)
		}
	},
}

func init() {
	rootCmd.AddCommand(showCmd)

	showCmd.Flags().IntVarP(&daysSince, "since-days", "s", 7, "show since input days back")
	showCmd.Flags().Bool("with-time", false, "prints with the time when the entry was added")
}
