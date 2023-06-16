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

var pyScripts = &core.Scripts{
	Init: `
`,
	Gen: `
`,
	Start: `
`,
	Stop: `
`,
}

func Python(parent string) *core.Language {
	return core.CreateLanguage(
		"python",
		parent,
		pyScripts,
		&pyTpl,
		pyTypes,
		"List[%s]",
	)
}
