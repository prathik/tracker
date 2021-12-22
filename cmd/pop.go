package cmd

import (
	"github.com/prathik/tracker/domain"
	"github.com/prathik/tracker/repo"

	"github.com/spf13/cobra"
)

// popCmd represents the pop command
var popCmd = &cobra.Command{
	Use:   "pop",
	Short: "Pops last added entry",
	Run: func(cmd *cobra.Command, args []string) {
		db := cmd.Flag("db").Value.String()
		bolt := repo.NewBoltDbRepo(db)
		defer bolt.Close()
		ss := domain.NewSessionService(bolt)
		ss.Pop()
	},
}

func init() {
	rootCmd.AddCommand(popCmd)
}
