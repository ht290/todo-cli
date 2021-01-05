package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var newUsername *string

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a notebook for the user",
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		password, err := getPassword()
		if err != nil {
			return err
		}
		return initNotebook().Create(*newUsername, password)
	},
}

func getPassword() (string, error) {
	fmt.Println("Enter password")
	var password string
	_, err := fmt.Scan(&password)
	if err != nil {
		return "", err
	}
	return password, nil
}

func init() {
	rootCmd.AddCommand(createCmd)
	newUsername = createCmd.Flags().StringP("user", "u", "", "Username")
}
