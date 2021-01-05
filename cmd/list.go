package cmd

import (
	"fmt"
	"todo-cli/notebook"

	"github.com/spf13/cobra"
)

var listAllItems *bool

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List todo items",
	RunE: func(cmd *cobra.Command, args []string) error {
		if *listAllItems {
			return listAll()
		}
		return listUndone()
	},
}

func listAll() error {
	items, doneItems, err := initNotebook().ListAllItems()
	if err != nil {
		return err
	}
	printItems(items)
	fmt.Printf("Total: %d items, %d items done\n", len(items), doneItems)
	return nil
}

func listUndone() error {
	items, err := initNotebook().ListUndoneItems()
	if err != nil {
		return err
	}
	printItems(items)
	fmt.Printf("Total: %d items\n", len(items))
	return nil
}

func printItems(items []notebook.Item) {
	for _, item := range items {
		doneLabel := " "
		if item.Done {
			doneLabel = " [Done] "
		}
		fmt.Printf("%v.%s%s\n", item.Id, doneLabel, item.Summary)
	}
}

func init() {
	rootCmd.AddCommand(listCmd)
	listAllItems = listCmd.Flags().BoolP("all", "a", false, "Include done items")
}
