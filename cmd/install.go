/*
Copyright © 2023 Jonathan Taylor <jonrtaylor12@gmail.com>
*/

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// installCmd represents the installation command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "install installs all module dependencies for a given system",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("install called")
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}
