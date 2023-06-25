/*
Copyright © 2023 Jonathan Taylor <jonrtaylor12@gmail.com>
*/

package cmd

import (
	"fmt"
	"github.com/jt05610/scaf/builder"
	"github.com/jt05610/scaf/context"
	"github.com/jt05610/scaf/core"
	"github.com/jt05610/scaf/lang"
	"github.com/jt05610/scaf/yaml"
	"github.com/spf13/cobra"
	uz "go.uber.org/zap"
	"os"
)

func Build(ctx context.Context, parent string, s *core.System) {
	ctx.Logger.Info("building system", uz.String("name", s.Name))
	goLang := lang.Go(parent)
	ts := lang.TypeScript(parent)

	bld := builder.NewBuilder(
		builder.NewRunner(parent, core.Go, goLang.CmdSet, builder.ModScriptBuild),
		builder.NewRunner(parent, core.TypeScript, ts.CmdSet, builder.ModScriptBuild),
	)

	err := bld.VisitSystem(ctx, s)
	if err != nil {
		ctx.Logger.Error("failed to build system", uz.Error(err))
	}
	ctx.Logger.Info("system built", uz.String("name", s.Name))

}

// genCmd represents the gen command
var buildCmd = &cobra.Command{
	Use:   "gen",
	Short: "generate code for a system",
	Long:  `Rather than write code for a system by hand, scaf can generate almost all of it for you!`,
	Run: func(cmd *cobra.Command, args []string) {
		if sysConfig == "" {
			fmt.Println("Please provide a system configuration file or run `scaf init` to create one.")
			return
		}
		ctx := Ctx(parDir, "gen")
		df, err := os.Open(sysConfig)
		if err != nil {
			ctx.Logger.Error("", uz.Error(err))
			return
		}
		srv := yaml.Service[*core.System]{}
		sys, err := srv.Load(df)
		if err != nil {
			ctx.Logger.Error("error loading system config", uz.Error(err))
			return
		}
		gen(ctx, parDir, sys)
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)
}
