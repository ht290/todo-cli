package cmd

import (
	"github.com/spf13/cobra"
	"todo-cli/notebook"
)

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Close the notebook for the current user",
	RunE: func(cmd *cobra.Command, args []string) error {
		return notebook.New().Lock()
	},
}

func init() {
	rootCmd.AddCommand(logoutCmd)
}
