/*
Copyright © 2023 Jonathan Taylor <jonrtaylor12@gmail.com>
*/

package cmd

import (
	"fmt"
	"github.com/jt05610/scaf/core"
	_go "github.com/jt05610/scaf/go"
	"github.com/jt05610/scaf/python"
	"github.com/spf13/cobra"
	uz "go.uber.org/zap"
	"os"
	"strings"
)

var modName, modDesc, modAuth, modLang string

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new module to the system",
	Long: `Add a new module to the system. This will create a new directory in the system's root 
directory and a new module configuration file.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := Ctx(parDir, "add")
		ctx.Logger.Debug("loading system")
		system, err := LoadSystem()
		if err != nil {
			ctx.Logger.Error("failed to load system", uz.Error(err))
			return
		}
		ctx.Logger.Debug("system loaded", uz.String("name", system.Name))

		if modName == "" {
			ctx.Logger.Debug("no module name provided, prompting user")
			err := getValue("Module name", &modName)
			if err != nil {
				ctx.Logger.Error("failed to get module name", uz.Error(err))
				return
			}
		}
		if modDesc == "" {
			ctx.Logger.Debug("no module description provided, prompting user")
			err := getValue("Module description", &modDesc)
			if err != nil {
				ctx.Logger.Error("failed to get module description", uz.Error(err))
				return
			}
		}
		if modAuth == "" {
			ctx.Logger.Debug("no module author provided, prompting user")
			defaultAuthor := DefaultAuthor(ctx)
			err := getValue(fmt.Sprintf("Author (default: %s)", defaultAuthor), &modAuth)
			if err != nil {
				ctx.Logger.Error("failed to get module author", uz.Error(err))
				return
			}
			if modAuth == "" {
				modAuth = defaultAuthor
			}
		}
		var lang core.Language
		if modLang == "" {
			ctx.Logger.Debug("no module language provided, prompting user")
			err := getValue("Module language", &modLang)
			if err != nil {
				ctx.Logger.Error("failed to get module language", uz.Error(err))
				return
			}
		}
		switch strings.ToLower(modLang) {
		case "go", "":
			lang = _go.Lang(parDir)
		case "python", "py":
			lang = python.Lang(parDir)
		default:
			ctx.Logger.Error("unsupported language", uz.String("lang", modLang))
		}

		ctx.Logger.Info("adding module", uz.String("name", modName))
		err = system.AddModule(core.NewModule(modName, modDesc, modAuth, lang))
		if err != nil {
			ctx.Logger.Error("failed to add module", uz.Error(err))
			return
		}
		ctx.Logger.Info("module added", uz.String("name", modName))

		ctx.Logger.Info("creating module directory")
		err = os.MkdirAll(fmt.Sprintf("%s/%s", parDir, modName), os.ModePerm)
		if err != nil {
			ctx.Logger.Error("failed to create module directory", uz.Error(err))
			return
		}
		ctx.Logger.Info("module directory created", uz.String("name", modName))
		ctx.Logger.Info("saving system")
		err = SaveSystem(system)
		if err != nil {
			ctx.Logger.Error("failed to save system", uz.Error(err))
			return
		}
		ctx.Logger.Info("system saved", uz.String("name", system.Name))
	},
}

func getValue(prompt string, v *string) error {
	if *v == "" {
		fmt.Printf(" %s: ", prompt)
		if _, err := fmt.Scanln(v); err != nil {
			return err
		}
	}
	return nil
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().StringVarP(&modName, "name", "n", "", "The name of the module")
	addCmd.Flags().StringVarP(&modDesc, "description", "d", "", "A description of the module")
	addCmd.Flags().StringVarP(&modAuth, "author", "a", "", "The author of the module")
	addCmd.Flags().StringVarP(&modLang, "language", "l", "", "The language of the module")
}
