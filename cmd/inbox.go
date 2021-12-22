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

		inbox := domain.NewInbox(bolt)
		for _, item := range inbox.CapturedItems {
			fmt.Println(item.Text)
		}
	},
}

func init() {
	rootCmd.AddCommand(inboxCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// inboxCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// inboxCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
