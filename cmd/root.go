/*
Copyright © 2023 Jonathan Taylor jonrtaylor12@gmail.com
*/

package cmd

import (
	"github.com/jt05610/scaf/context"
	"github.com/jt05610/scaf/zap"
	"os"

	"github.com/spf13/cobra"
)

var cfgFile string
var debug bool
var parDir string

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
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", true, "verbose output")
}

func Ctx(call string) context.Context {
	var log *zap.Logger
	if debug {
		log = zap.NewDev(context.Background(), parDir, call)
	} else {
		log = zap.NewProd(context.Background(), parDir, call)
	}
	return context.NewContext(log)
}
