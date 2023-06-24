package lang

import (
	"embed"
	"github.com/jt05610/scaf/core"
)

//go:embed template/fortran
var fortranTpl embed.FS

var fortranTypes = TypeMap{
	core.Int:    "integer",
	core.Float:  "real",
	core.String: "character(len=20)",
	core.ID:     "character(len=20)",
	core.Bool:   "logical",
}

var fortranScripts = &Scripts{
	Map: map[string]string{
		"start": `
{{.Name}} serve --port {{.Port}} 
`,
		"stop": `
kill $(lsof -t -i:{{.Port}})
`,
	},
}

func Fortran(parent string) *Language {
	return CreateLanguage(
		"fortran",
		parent,
		nil,
		fortranScripts,
		&fortranTpl,
		fortranTypes,
		"type(%s), allocatable, dimension(:)",
	)
}
