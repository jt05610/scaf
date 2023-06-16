/*
Copyright © 2023 Jonathan Taylor <jonrtaylor12@gmail.com>
*/

package cmd

import (
	"github.com/jt05610/scaf/context"
	"github.com/jt05610/scaf/core"
	"github.com/jt05610/scaf/yaml"
	"go.uber.org/zap"
	"os"
	"os/user"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)

func Init(ctx context.Context, name string) {
	ctx.Logger.Info("initializing system", zap.String("name", name))
	u, err := user.Current()
	if err != nil {
		ctx.Logger.Error("failed to get current user", zap.Error(err))
	}

	sys := &core.System{
		Name:    name,
		UIPort:  4000,
		GQLPort: 8080,
		Author:  u.Name,
		Date:    time.Now().Format("02 Jan 2006"),
		Modules: nil,
	}
	path := filepath.Join(parDir, name)
	err = os.MkdirAll(path, 0755)
	if err != nil {
		ctx.Logger.Error("failed to create directory", zap.String("path", path), zap.Error(err))
	}
	sysFile := filepath.Join(path, "system.yaml")
	f, err := os.Create(sysFile)
	if err != nil {
		ctx.Logger.Error("failed to create file", zap.String("path", sysFile), zap.Error(err))
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			ctx.Logger.Error("failed to close file", zap.String("path", sysFile), zap.Error(err))
		}
	}(f)
	srv := yaml.NewYAMLService[*core.System]()
	err = srv.Flush(f, sys)
	if err != nil {
		ctx.Logger.Error("failed to flush system", zap.String("path", sysFile), zap.Error(err))
	}
	ctx.Logger.Info("system initialized", zap.String("path", sysFile))
}

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := Ctx("init")
		Init(ctx, args[0])
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

}
