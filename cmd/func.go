/*
Copyright © 2023 Jonathan Taylor <jonrtaylor12@gmail.com>
*/

package cmd

import (
	"github.com/spf13/cobra"
)

// modelCmd represents the model command
var funcCmd = &cobra.Command{
	Use:   "func",
	Short: "add a function to a module",
	Long:  `Func is used to add a new function to a module.`,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	funcCmd.AddCommand(modelCmd)
}
