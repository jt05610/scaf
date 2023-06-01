/*
Copyright © 2023 Jonathan Taylor jonrtaylor12@gmail.com
*/

package cmd

import (
	"github.com/jt0610/scaf/context"
	"github.com/jt0610/scaf/wizard"
	"github.com/jt0610/scaf/zap"
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
		ctx = context.NewContext(zap.NewProd(context.Background(), "init"))
		w := &wizard.Wizard{}
		err := w.Run(ctx, nil)
		if err != nil {
			ctx.Logger.Panic("Error running init: " + err.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
