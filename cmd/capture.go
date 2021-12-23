package cmd

import (
	"github.com/fatih/color"
	"github.com/prathik/tracker/domain"
	"github.com/prathik/tracker/repo"
	"time"

	"github.com/spf13/cobra"
	"github.com/manifoldco/promptui"
)

// captureCmd represents the capture command
var captureCmd = &cobra.Command{
	Use:   "capture",
	Short: "Capture an inbox item",
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		db := cmd.Flag("db").Value.String()
		bolt, err := repo.NewBoltDbRepo(db)
		if err != nil {
			color.Red("error: %s", err)
			return
		}
		defer bolt.Close()
		prompt := promptui.Prompt{
			Label:    "Inbox",
		}

		item, err := prompt.Run()
		inboxItem := domain.NewInboxItem(time.Now(), item, bolt)
		err = inboxItem.Save()
		if err != nil {
			color.Red("error: %s", err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(captureCmd)
}
