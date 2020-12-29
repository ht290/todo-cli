/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
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
	userNotebook := notebook.NewForSingleUser()
	items, doneItems, err := userNotebook.ListAllItems()
	if err != nil {
		return err
	}
	printItems(items)
	fmt.Printf("Total: %d items, %d items done\n", len(items), doneItems)
	return nil
}

func listUndone() error {
	userNotebook := notebook.NewForSingleUser()
	items, err := userNotebook.ListUndoneItems()
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
