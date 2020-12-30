package cmd

import (
	"fmt"
	"strconv"
	"todo-cli/notebook"

	"github.com/spf13/cobra"
)

var doneCmd = &cobra.Command{
	Use:   "done",
	Short: "Complete a todo item",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		notebook := notebook.New()
		// No need to check len(args) since `cobra.ExactArgs(1)` guarantees exactly 1 arg
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}
		if err := notebook.Done(id); err != nil {
			return err
		}
		fmt.Printf("Item %v done.\n", id)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(doneCmd)
}
