package cmd

import (
	"github.com/spf13/cobra"
	"todo-cli/notebook"
	"todo-cli/notebook/persistence"
)

var importFilename *string
var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import todo items from a Json file",
	RunE: func(cmd *cobra.Command, args []string) error {
		itemFile := persistence.JsonFile{
			Filename: *importFilename,
		}
		var items []notebook.Item
		if err := itemFile.Read(&items); err != nil {
			return err
		}
		myNotebook := initNotebook()
		for _, item := range items {
			if _, err := myNotebook.Add(item.Summary); err != nil {
				return err
			}
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(importCmd)
	importFilename = importCmd.Flags().StringP("filename", "f", "", "Json file to import")
}
