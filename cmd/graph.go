/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/prathik/tracker/repo"
	"github.com/prathik/tracker/domain"
	"time"

	"github.com/guptarohit/asciigraph"
	"github.com/spf13/cobra"
)

// graphCmd represents the graph command
var graphCmd = &cobra.Command{
	Use:   "graph",
	Short: "Plots the graph for importance and joy",
	Run: func(cmd *cobra.Command, args []string) {
		db := cmd.Flag("db").Value.String()
		bolt := repo.NewBoltDbRepo(db)
		defer bolt.Close()
		ss := domain.NewSessionService(bolt)
		duration, _ := time.ParseDuration("168h") // 7 days
		queryData := ss.QueryData(duration)
		var importance []float64
		var joy []float64
		for _, v := range queryData.Days {
			for _, d := range v.Sessions {
				importance = append(importance, float64(d.Impact))
				joy = append(joy, float64(d.Joy))
			}
		}
		impGraph := asciigraph.Plot(importance)

		color.Green("Impact")
		fmt.Println(impGraph)

		joyGraph := asciigraph.Plot(joy)
		color.Green("\nJoy")
		fmt.Println(joyGraph)

	},
}

func init() {
	showCmd.AddCommand(graphCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// graphCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// graphCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
