package lang

import (
	"embed"
	"github.com/jt05610/scaf/core"
)

//go:embed template/python
var pyTpl embed.FS

var pyTypes = &core.TypeMap{
	Int:    "int",
	Float:  "float",
	String: "str",
	Bool:   "bool",
}

var shService = `
	srvCmd := exec.Command("/bin/sh", "scripts/service.sh")

	srvCmd.Stdout = os.Stdout
	srvCmd.Stderr = os.Stderr

	if err := srvCmd.Run(); err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
`

var pyScripts = &core.Scripts{
	Init: `
sh scripts/init.sh
`,
	Gen: `
sh scripts/gen.sh
`,
}

func Python(parent string) *core.Language {
	l := core.CreateLanguage(
		"python",
		parent,
		pyScripts,
		&pyTpl,
		pyTypes,
		"List[%s]",
	)
	l.Service = shService
	return l
}
