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
	"encoding/json"
	"fmt"
	"todo-cli/notebook"

	"github.com/spf13/cobra"
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export todo items to Json",
	RunE: func(cmd *cobra.Command, args []string) error {
		allItems, _, err := initNotebook().ListAllItems()
		if err != nil {
			return err
		}
		return printJson(allItems)
	},
}

func printJson(items []notebook.Item) error {
	jsonItems, err := json.MarshalIndent(items, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(jsonItems))
	return nil
}

func init() {
	rootCmd.AddCommand(exportCmd)
}
