package lang

import (
	"embed"
	"github.com/jt05610/scaf/core"
)

//go:embed template/python
var pyTpl embed.FS

var pyTypes = TypeMap{
	core.Int:    "int",
	core.Float:  "float",
	core.String: "str",
	core.ID:     "str",
	core.Bool:   "bool",
}

var shService = `
	srvCmd := exec.Command("/bin/sh", "scripts/service.sh")

	srvCmd.Stdout = os.Stdout
	srvCmd.Stderr = os.Stderr

	if err := srvCmd.Run(); err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
`

var pyScripts = &Scripts{
	Map: map[string]string{
		"init": `
sh scripts/init.sh
`,
		"gen": `
sh scripts/gen.sh
`,
	},
}

func Python(parent string) *Language {
	l := CreateLanguage(
		"python",
		parent,
		nil,
		pyScripts,
		&pyTpl,
		pyTypes,
		"List[%s]",
	)
	l.Service = shService
	return l
}
