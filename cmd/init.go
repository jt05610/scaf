/*
Copyright © 2023 Jonathan Taylor jonrtaylor12@gmail.com
*/

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes a new robot with appropriate folder structure.",
	Long: `The 'init' command of scaf initializes a new robot by creating the required folder 
structure and config files based on the robot description provided. This structure allows for efficient organization 
and operation of the diverse modules and services associated with the robot.

Running 'scaf init' will:

- Download the core libraries and modules required for your robot.
- Create subdirectories for each module described in the YAML file.
- Set up a GraphQL file for defining API endpoints.
- Initiate spaces for firmware, SolidJS microfrontend, and microservices related to the robot.

By default, scaf generates blinky: a robot with the ability to control the color, brightness, and 
frequency of the blinking LED. This can be switched off with the --no-blinky flag.`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("init called")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
