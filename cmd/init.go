/*
Copyright © 2023 Jonathan Taylor jonrtaylor12@gmail.com
*/

package cmd

import (
	"github.com/jt05610/scaf/context"
	"github.com/jt05610/scaf/system"
	sz "github.com/jt05610/scaf/zap"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"os"
	"path/filepath"
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
		ctx = context.NewContext(sz.NewProd(context.Background(), parDir, "init"))
		err := os.MkdirAll(parDir, 0755)
		if err != nil {
			ctx.Logger.Panic("could not create parent directory", zap.Error(err))
		}
		f, err := os.Create(filepath.Join(parDir, "system.yaml"))
		if err != nil {
			ctx.Logger.Panic("could not create system.yaml", zap.Error(err))
		}
		defer func(f *os.File) {
			err := f.Close()
			if err != nil {
				ctx.Logger.Panic("could not close system.yaml", zap.Error(err))
			}
		}(f)
		sys, err := hndl.CreateSystem(ctx, f)
		if err != nil {
			ctx.Logger.Panic("could not create system", zap.Error(err))
		}
		ctx.Logger.Info("created system", zap.String("name", sys.Name))
		ctx.Logger.Info("wrote system.yaml")
		err = os.MkdirAll(filepath.Join(parDir, "config"), 0755)
		if err != nil {
			ctx.Logger.Panic("could not create config directory", zap.Error(err))
		}
		ctx.Logger.Info("created config directory")
		if sys.Kind != system.Library {
			cf := sys.Caddyfile()
			caddyConfig, err := os.Create(filepath.Join(parDir, "config/Caddyfile.yaml"))

			if err != nil {
				ctx.Logger.Panic("could not create Caddyfile config", zap.Error(err))
			}
			defer func(f *os.File) {
				err := f.Close()
				if err != nil {
					ctx.Logger.Panic("could not close Caddyfile config", zap.Error(err))
				}
			}(caddyConfig)
			err = hndl.CaddyService.Flush(caddyConfig, cf)
			if err != nil {
				ctx.Logger.Panic("could not flush Caddyfile config", zap.Error(err))
			}
			ctx.Logger.Info("created config/Caddyfile.yaml")

			nf, err := os.Create(filepath.Join(parDir, "Caddyfile"))
			if err != nil {
				ctx.Logger.Panic("could not create Caddyfile", zap.Error(err))
			}
			defer func(f *os.File) {
				err := f.Close()
				if err != nil {
					ctx.Logger.Panic("could not close Caddyfile", zap.Error(err))
				}
			}(nf)
			err = hndl.CaddyRenderer.Flush(nf, cf)
			if err != nil {
				ctx.Logger.Panic("could not flush Caddyfile", zap.Error(err))
			}
			ctx.Logger.Info("created Caddyfile")
		}

	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
