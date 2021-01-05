package cmd

import (
	"github.com/spf13/cobra"
)

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Close the notebook for the current user",
	RunE: func(cmd *cobra.Command, args []string) error {
		return initNotebook().Lock()
	},
}

func init() {
	rootCmd.AddCommand(logoutCmd)
}
