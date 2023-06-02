/*
Copyright © 2023 Jonathan Taylor jonrtaylor12@gmail.com
*/

package cmd

import (
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"os"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes a new system with appropriate folder structure.",
	Long: `The 'init' command of scaf initializes a new system by creating the required folder 
structure and config files based on the system description provided. This structure allows for efficient organization 
and operation of the diverse modules and services associated with the system.

Running 'scaf init' will:

- Download the core libraries and modules required for your system.
- Create subdirectories for each module described in the YAML file.
- Set up a GraphQL file for defining API endpoints.
- Initiate spaces for firmware, SolidJS microfrontend, and microservices related to the system.

By default, scaf generates blinky: a system with the ability to control the color, brightness, and 
frequency of the blinking LED. This can be switched off with the --no-blinky flag.`,

	Run: func(cmd *cobra.Command, args []string) {
		f, err := os.Create("system.yaml")
		if err != nil {
			ctx.Logger.Panic("could not create system.yaml", zap.Error(err))
		}
		defer func(f *os.File) {
			err := f.Close()
			if err != nil {
				ctx.Logger.Panic("could not close system.yaml", zap.Error(err))
			}
		}(f)
		err = hndl.CreateSystem(ctx, f)
		if err != nil {
			ctx.Logger.Panic("could not create system", zap.Error(err))
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
