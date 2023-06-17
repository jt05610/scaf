/*
Copyright © 2023 Jonathan Taylor <jonrtaylor12@gmail.com>
*/

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var noCrud bool

// modelCmd represents the model command
var modelCmd = &cobra.Command{
	Use:   "model",
	Short: "add a model to a module",
	Long: `Model is used to add a new model to a module. By default, modules will implement all CRUD operations for 
models. This can be overridden by using the --no-crud flag.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("model called")
	},
}

func init() {
	addCmd.AddCommand(modelCmd)
	modelCmd.Flags().BoolVar(&noCrud, "no-crud", false, "do not generate CRUD operations for this model")
}
