package cmd

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/prathik/tracker/domain"
	"github.com/prathik/tracker/repo"
	"github.com/spf13/cobra"
)

// inboxCmd represents the inbox command
var inboxCmd = &cobra.Command{
	Use:   "inbox",
	Short: "Show inbox",
	Run: func(cmd *cobra.Command, args []string) {
		db := cmd.Flag("db").Value.String()
		bolt, err := repo.NewBoltDbRepo(db)
		if err != nil {
			color.Red("error: %s", err)
			return
		}
		defer bolt.Close()

		inbox := domain.NewInbox(bolt)
		for _, item := range inbox.CapturedItems {
			fmt.Println(item.Text)
		}
	},
}

func init() {
	rootCmd.AddCommand(inboxCmd)
}
