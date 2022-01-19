package cmd

import (
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/prathik/tracker/domain"
	"github.com/prathik/tracker/repo"
	"os"
	"time"

	"github.com/spf13/cobra"
)

// statsCmd represents the stats command
var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		db := cmd.Flag("db").Value.String()
		bolt, err := repo.NewBoltDbRepo(db)
		if err != nil {
			color.Red("error: %s", err)
			return
		}
		defer bolt.Close()
		sessionService := domain.NewSessionService(bolt)
		PrintWeeklyReport(sessionService)
	},
}

func PrintWeeklyReport(sessionService *domain.SessionService) {
	data, _ := domain.GenerateWeekReport(time.Now(), sessionService)
	color.Green("Average Count = %.2f", data.AverageCount)
	color.Green("Flow Ratio = %.2f", data.FlowRatio)
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Day", "Count"})
	table.AppendBulk(data.Raw)
	table.Render()
}

func init() {
	showCmd.AddCommand(statsCmd)
}
