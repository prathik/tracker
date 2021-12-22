package cmd

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/prathik/tracker/domain"
	"github.com/prathik/tracker/repo"
	"time"

	"github.com/guptarohit/asciigraph"
	"github.com/spf13/cobra"
)

// graphCmd represents the graph command
var graphCmd = &cobra.Command{
	Use:   "graph",
	Short: "Plots the graph for flow",
	Run: func(cmd *cobra.Command, args []string) {
		db := cmd.Flag("db").Value.String()
		bolt, err := repo.NewBoltDbRepo(db)
		if err != nil {
			color.Red("error: %s", err)
			return
		}

		defer bolt.Close()
		ss := domain.NewSessionService(bolt)
		duration, _ := time.ParseDuration("168h") // 7 days
		queryData, err := ss.QueryData(duration)
		if err != nil {
			color.Red("error: %s", err)
			return
		}
		var flow []float64
		for _, v := range queryData.Days {
			for _, session := range v.Sessions {
				flow = append(flow, float64(domain.Score(session.Challenge)))
			}
		}
		flowGraph := asciigraph.Plot(flow)

		color.Green("Flow")
		fmt.Println(flowGraph)
	},
}

func init() {
	showCmd.AddCommand(graphCmd)
}
