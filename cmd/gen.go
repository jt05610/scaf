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
	_go "github.com/jt05610/scaf/go"
	"github.com/jt05610/scaf/gql"
	"github.com/jt05610/scaf/proto"
	"github.com/jt05610/scaf/ts"
	"github.com/jt05610/scaf/yaml"
	"github.com/spf13/cobra"
	uz "go.uber.org/zap"
	"os"
)

var sysConfig string

func Gen(ctx context.Context, parent string, s *core.System) {
	ctx.Logger.Info("generating system", uz.String("name", s.Name))
	goLang := _go.Lang(parent)
	GQL := gql.Lang(parent)
	protoBuf := proto.Lang(parent)
	typescript := ts.Lang(parent)
	bld := builder.NewBuilder(
		codegen.New(parent, goLang),
		codegen.New(parent, GQL),
		codegen.New(parent, protoBuf),
		codegen.New(parent, typescript),
		builder.NewRunner(parent, protoBuf.Gen()),
		builder.NewRunner(parent, goLang.Gen()),
		builder.NewRunner(parent, GQL.Gen()),
		builder.NewRunner(parent, typescript.Gen()),
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
