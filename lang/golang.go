package lang

import (
	"embed"
	"github.com/jt05610/scaf/core"
)

//go:embed template/go
var goTpl embed.FS

var goTypes = &core.TypeMap{
	Int:    "int64",
	Float:  "float64",
	String: "string",
	Bool:   "bool",
}

var goScripts = &core.Scripts{
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

func Go(parent string) *core.Language {
	return core.CreateLanguage(
		"go",
		parent,
		goScripts,
		&goTpl,
		goTypes,
		"[]%s",
	)
}
