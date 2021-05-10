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
	"errors"
	"fmt"
	"github.com/manifoldco/promptui"
	"github.com/olekukonko/tablewriter"
	"github.com/prathik/tracker/repo"
	"github.com/prathik/tracker/service"
	"github.com/spf13/cobra"
	"os"
	"strconv"
	"time"
)

// incrementCmd represents the increment command
var incrementCmd = &cobra.Command{
	Use:   "increment",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		bolt := repo.NewBoltDbRepo()
		defer bolt.Close()
		ss := service.NewSessionService(bolt)

		workPrompt := promptui.Prompt{
			Label: "Work",
			Validate: func(s string) error {
				if len(s) < 7 {
					return errors.New("enter a value")
				}

				return nil
			},
		}
		workResult, _ := workPrompt.Run()

		joy := loadInteger("Joy [0-10]")

		imp := loadInteger("Importance [0-10]")

		notesPrompt := promptui.Prompt{
			Label:   "Notes",
			Default: "",
		}
		notesResult, _ := notesPrompt.Run()

		ss.Create(&service.Item{Work: workResult, Joy: joy, Importance: imp, Notes: notesResult, Time: time.Now()})

		fmt.Printf("\n")

		// Print entire week as a table for a reminder
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Day", "Count"})

		for _, d := range ss.GetWeekData().WeekData {
			table.Append([]string{d.Time.Format("2006-01-02"), strconv.Itoa(d.Count)})
		}
		table.Render()
	},
}

func loadInteger(label string) int {
	intPrompt := promptui.Prompt{
		Label:    label,
		Validate: validateNumber,
	}
	result, _ := intPrompt.Run()
	intVal, _ := strconv.Atoi(result)
	return intVal
}

func validateNumber(input string) error {
	ip, err := strconv.Atoi(input)
	if err != nil {
		return errors.New("invalid number")
	}

	if ip > 10 || ip < 0 {
		return errors.New("enter between 0 to 10")
	}
	return nil
}

func init() {
	rootCmd.AddCommand(incrementCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// incrementCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// incrementCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
