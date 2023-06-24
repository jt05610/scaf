/*
Copyright © 2023 Jonathan Taylor jonrtaylor12@gmail.com
*/

package cmd

import (
	"fmt"
	"github.com/jt05610/scaf/context"
	"github.com/jt05610/scaf/core"
	"github.com/jt05610/scaf/yaml"
	"github.com/jt05610/scaf/zap"
	uz "go.uber.org/zap"
	"os"
	"os/user"
	"path/filepath"

	"github.com/spf13/cobra"
)

var cfgFile string
var debug bool
var parDir string

func FakeSystems(n int) []*core.System {
	ret := make([]*core.System, n)
	for i := 0; i < n; i++ {
		ret[i] = core.NewSystem(fmt.Sprintf("System %d", i), "", "Jonathan Taylor", fmt.Sprintf("%d Jul 2023", i%31))
	}
	return ret
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "scaf",
	Short: "scaf is a Systems Configuration Automation Framework for robotics systems.",
	Long: `scaf, standing for Systems Configuration Automation Framework, is a revolutionary meta-robotics tool 
that automates the construction and management of complex robotics systems. It efficiently integrates various 
hardware and software modules, using YAML and GraphQL descriptions for automatic firmware generation, microfrontends, 
and microservices. It's ideal for researchers and developers who want to streamline their robotics and automation processes.

For example, with scaf, you can:

- Describe your robotics modules in YAML files and define API endpoints with GraphQL.
- Generate firmware automatically based on the provided system descriptions.
- Create SolidJS microfrontends that can be easily incorporated into a primary UI app frame.
- Produce microservices for the new modules in your chosen language, with support for Python and Go, and extensibility for other languages.
- Download and organize your required modules automatically in a local directory, akin to Python's venv or Node's node_modules.

`,
}

func Execute() {
	err := rootCmd.Execute()

	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "$HOME/.scaf.yaml", "config file (default is $HOME/.scaf.yaml)")
	rootCmd.PersistentFlags().StringVarP(&parDir, "parent", "p", ".", "parent directory for system")
	rootCmd.PersistentFlags().BoolVarP(&debug, "verbose", "v", false, "verbose output")
}

func Ctx(parent string, call string) context.Context {
	var log *zap.Logger
	if debug {
		log = zap.NewDev(context.Background(), parent, call)
	} else {
		log = zap.NewProd(context.Background(), parent, call)
	}
	return context.NewContext(log)
}

func DefaultAuthor(ctx context.Context) string {
	u, err := user.Current()
	if err != nil {
		ctx.Logger.Error("failed to get current user", uz.Error(err))
	}
	return u.Name
}

func LoadSystem() (*core.System, error) {
	path := filepath.Join(parDir, "system.yaml")
	fp, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = fp.Close()
	}()
	srv := yaml.NewYAMLService[*core.System]()
	return srv.Load(fp)
}

func SaveSystem(s *core.System) error {
	path := filepath.Join(parDir, "system.yaml")
	fp, err := os.Create(path)
	if err != nil {
		return err
	}
	srv := yaml.NewYAMLService[*core.System]()
	err = srv.Flush(fp, s)
	if err != nil {
		return err
	}
	return fp.Close()
}
