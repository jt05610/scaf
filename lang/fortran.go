package lang

import (
	"embed"
	"github.com/jt05610/scaf/core"
)

//go:embed template/fortran
var fortranTpl embed.FS

var fortranTypes = &core.TypeMap{
	Int:    "integer",
	Float:  "real",
	String: "character(len=20)",
	Bool:   "logical",
}

var fortranScripts = &core.Scripts{
	Init: `
`,
	Gen: `
mkcert {{.Name}}.local localhost 127.0.0.1 ::1
mkdir ./cmd/.secrets
mv {{.Name}}.local+3.pem {{.Name}}.local+3-key.pem ./cmd/.secrets
`,
	Start: `
{{.Name}} serve --port {{.Port}} 
`,
	Stop: `
kill $(lsof -t -i:{{.Port}})
`,
}

func Fortran(parent string) *core.Language {
	return core.CreateLanguage(
		"fortran",
		parent,
		nil,
		fortranScripts,
		&fortranTpl,
		fortranTypes,
		"type(%s), allocatable, dimension(:)",
	)
}
