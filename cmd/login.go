package cmd

import (
	"github.com/spf13/cobra"
	"todo-cli/notebook"
)

var username *string

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Open the notebook for the user",
	RunE: func(cmd *cobra.Command, args []string) error {
		password, err := getPassword()
		if err != nil {
			return err
		}
		return notebook.New().Unlock(*username, password)
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
	username = loginCmd.Flags().StringP("user", "u", "", "Username")
}
