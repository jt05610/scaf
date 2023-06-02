/*
Copyright © 2023 Jonathan Taylor <jonrtaylor12@gmail.com>
*/

package cmd

import (
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"os"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "add a new module to the system",
	Long: `
The 'add' command of scaf adds a new module to the system. This module can hardware or software.
`,
	Run: func(cmd *cobra.Command, args []string) {
		sf, err := os.Open("system.yaml")
		if err != nil {
			ctx.Logger.Panic("could not open system.yaml", zap.Error(err))
		}
		sys, err := hndl.SystemService.Load(sf)
		if err != nil {
			ctx.Logger.Panic("could not load system", zap.Error(err))
		}
		ctx.Logger.Info("loaded system", zap.String("name", sys.Name))
		err = sf.Close()
		if err != nil {
			ctx.Logger.Panic("could not close system.yaml", zap.Error(err))
		}
		m, err := hndl.CreateModule(ctx)
		if err != nil {
			ctx.Logger.Panic("could not create module", zap.Error(err))
		}
		sys.AddModule(m)
		sf, err = os.Create("system.yaml")
		if err != nil {
			ctx.Logger.Panic("could not flush system.yaml", zap.Error(err))
		}
		err = hndl.SystemService.Flush(sf, sys)
		if err != nil {
			ctx.Logger.Panic("could not flush system.yaml", zap.Error(err))
		}
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
