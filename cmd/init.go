/*
Copyright © 2023 Jonathan Taylor <jonrtaylor12@gmail.com>
*/

package cmd

import (
	"fmt"
	"github.com/jt05610/scaf/context"
	"github.com/jt05610/scaf/core"
	"github.com/jt05610/scaf/yaml"
	"go.uber.org/zap"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)

var name string

func Init(ctx context.Context, name string) {
	ctx.Logger.Info("initializing system", zap.String("name", name))

	sys := core.NewSystem(name, DefaultAuthor(ctx), time.Now().Format("02 Jan 2006"))
	path := filepath.Join(parDir, name)
	err := os.MkdirAll(path, 0755)
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
	Short: "initialize a new system",
	Long:  `Initialize a new system by creating a directory and system.yaml which will store the device definitions.`,
	Run: func(cmd *cobra.Command, args []string) {
		if name == "" {
			fmt.Print("Please provide a name for the system: ")
			if _, err := fmt.Scanln(&name); err != nil {
				fmt.Printf("Failed to read user input: %v\n", err)
				return
			}
		}
		path := filepath.Join(parDir, name)
		ctx := Ctx(path, "init")
		Init(ctx, name)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.PersistentFlags().StringVar(&name, "name", "", "name of the system")
}
