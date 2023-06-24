/*
Copyright © 2023 Jonathan Taylor <jonrtaylor12@gmail.com>
*/

package cmd

import (
	"fmt"
	"github.com/jt05610/scaf/builder"
	"github.com/jt05610/scaf/codegen"
	"github.com/jt05610/scaf/context"
	"github.com/jt05610/scaf/core"
	"github.com/jt05610/scaf/lang"
	"github.com/jt05610/scaf/yaml"
	"github.com/spf13/cobra"
	uz "go.uber.org/zap"
	"os"
)

var sysConfig string

func Gen(ctx context.Context, parent string, s *core.System) {
	if sysConfig == "" {
		sysConfig = parent + "/system.yaml"
	}
	ctx.Logger.Info("generating system", uz.String("name", s.Name))
	goLang := lang.Go(parent)
	gql := lang.GraphQL(parent)
	sql := lang.SQL(parent)
	proto := lang.Protobuf(parent)
	ts := lang.TypeScript(parent)
	bld := builder.NewBuilder(
		codegen.New(parent, goLang),
		codegen.New(parent, proto),
		codegen.New(parent, sql),
		codegen.New(parent, gql),
		codegen.New(parent, ts),
		//builder.NewRunner(parent, proto.CmdSet, builder.ModScriptInit, builder.ModScriptGen),
		//builder.NewRunner(parent, goLang.CmdSet, builder.ModScriptInit, builder.ModScriptGen),
		//builder.NewRunner(parent, ts.CmdSet, builder.ModScriptInit, builder.ModScriptGen),
	)

	err := bld.VisitSystem(ctx, s)
	if err != nil {
		ctx.Logger.Error("failed to generate system", uz.Error(err))
	}
	ctx.Logger.Info("system generated", uz.String("name", s.Name))

}

// genCmd represents the gen command
var genCmd = &cobra.Command{
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
		Gen(ctx, parDir, sys)
	},
}

func init() {
	rootCmd.AddCommand(genCmd)

	genCmd.Flags().StringVarP(&sysConfig, "config", "s", "", "system configuration file")
}
