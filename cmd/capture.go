package cmd

import (
	"github.com/fatih/color"
	"github.com/prathik/tracker/domain"
	"github.com/prathik/tracker/repo"
	"time"

	"github.com/spf13/cobra"
)

// captureCmd represents the capture command
var captureCmd = &cobra.Command{
	Use:   "capture",
	Short: "Capture an inbox item",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		db := cmd.Flag("db").Value.String()
		bolt, err := repo.NewBoltDbRepo(db)
		if err != nil {
			color.Red("error: %s", err)
			return
		}
		defer bolt.Close()

		inboxItem := domain.NewInboxItem(time.Now(), args[0], bolt)
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
