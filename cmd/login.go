package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var username *string

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Open the notebook for the user",
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		password, err := getPassword()
		if err != nil {
			return err
		}
		if err := initNotebook().Unlock(*username, password); err != nil {
			return err
		}
		fmt.Println("Login success!")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
	username = loginCmd.Flags().StringP("user", "u", "", "Username")
}
