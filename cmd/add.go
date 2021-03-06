package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a todo item",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// No need to check len(args) since `cobra.ExactArgs(1)` guarantees exactly 1 arg
		summary := args[0]
		newItem, err := initNotebook().Add(summary)
		if err != nil {
			return err
		}
		fmt.Printf("%v. %s\n", newItem.Id, newItem.Summary)
		fmt.Printf("Item %v added\n", newItem.Id)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
