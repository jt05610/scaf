/*
Copyright © 2023 Jonathan Taylor <jonrtaylor12@gmail.com>
*/

package cmd

import (
	"github.com/jt05610/scaf/lang/ts"
	"github.com/jt05610/scaf/system"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"os"
	"path/filepath"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "add a new module to the system",
	Long: `
The 'add' command of scaf adds a new module to the system. This module can hardware or software.
`,
	Run: func(cmd *cobra.Command, args []string) {
		sysPath := filepath.Join(parDir, "system.yaml")
		sf, err := os.Open(sysPath)
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
		modDir := filepath.Join(parDir, "modules", m.Name)
		err = os.MkdirAll(modDir, 0755)
		if err != nil {
			ctx.Logger.Panic("could not create module directory", zap.Error(err))
		}
		switch m.Language {
		case system.TypeScript:
			gen := ts.NewGenerator(modDir)
			err = gen.Generate(m)
		default:
			ctx.Logger.Panic("unsupported language", zap.String("language", string(m.Language)))
		}
		err = sys.AddModule(m)
		if err != nil {
			ctx.Logger.Panic("could not add module to system", zap.Error(err))
		}
		sf, err = os.Create(sysPath)
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
