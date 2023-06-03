package goLang

import "text/template"

func main() *template.Template {
	t, err := template.New("main.go").Parse(`
package main

import "{{.Name}}/cmd"

func main() {
	cmd.Execute()
}

`)
	if err != nil {
		panic(err)
	}
	return t
}

func goMod() *template.Template {
	t, err := template.New("go.mod").Parse(`
module {{.Name}}

go 1.20

require (
	github.com/google/go-cmp v0.5.9
	github.com/spf13/cobra v1.7.0
)
`)
	if err != nil {
		panic(err)
	}
	return t
}

func root() *template.Template {
	t, err := template.New("root.go").Parse(`
package cmd

import (
	"os"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "{{.Name}}",
	Short: "{{.Desc}}",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
}`)
	if err != nil {
		panic(err)

	}
	return t
}
