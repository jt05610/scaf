/*
Copyright © 2023 Jonathan Taylor <jonrtaylor12@gmail.com>
*/

package cmd

import (
	"github.com/jt05610/scaf/context"
	"github.com/jt05610/scaf/lang"
	"github.com/jt05610/scaf/lang/goLang"
	"github.com/jt05610/scaf/lang/py"
	"github.com/jt05610/scaf/lang/ts"
	"github.com/jt05610/scaf/system"
	sz "github.com/jt05610/scaf/zap"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"os"
	"path/filepath"
	"text/template"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "add a new module to the system",
	Long: `
The 'add' command of scaf adds a new module to the system. This module can hardware or software.
`,
	Run: func(cmd *cobra.Command, args []string) {

		ctx = context.NewContext(sz.NewProd(context.Background(), parDir, "add"))
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
		err = sys.AddModule(m)
		if err != nil {
			ctx.Logger.Panic("could not add module to system", zap.Error(err))
		}
		var files func(string) []*lang.File
		var install func() *template.Template
		switch m.Language {
		case system.TypeScript:
			files = ts.Files
		case system.Go:
			files = goLang.Files
		case system.Python:
			ctx.Logger.Info("creating python module")
			files = py.Files
			install = py.Install
		default:
			ctx.Logger.Panic("unsupported language", zap.String("language", string(m.Language)))
		}
		err = lang.Install(install(), &lang.Installable{ParDir: parDir})
		if err != nil {
			ctx.Logger.Panic("could not install module", zap.Error(err))
		}
		err = lang.Generate(modDir, files, m)
		if err != nil {
			ctx.Logger.Panic("could not generate module", zap.Error(err))
		}
		sf, err = os.Create(sysPath)
		if err != nil {
			ctx.Logger.Panic("could not flush system.yaml", zap.Error(err))
		}
		err = hndl.SystemService.Flush(sf, sys)
		if err != nil {
			ctx.Logger.Panic("could not flush system.yaml", zap.Error(err))
		}
		ctx.Logger.Info("created module", zap.String("name", m.Name))
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
