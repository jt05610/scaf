package lang

import (
	"embed"
	"github.com/jt05610/scaf/core"
)

//go:embed template/gql
var gqlTpl embed.FS

var gqlTypes = &core.TypeMap{
	Int:    "Int",
	Float:  "Float",
	String: "String",
	Bool:   "Boolean",
}

var gqlScripts = &core.Scripts{
	Init: `
go mod tidy
`,
	Gen: `
go generate ./...
go fmt ./...
`,
	Start: `
{{.Name}} serve --port {{.Port}} 
`,
	Stop: `
kill $(lsof -t -i:{{.Port}})
`,
}

func GraphQL(parent string) *core.Language {
	return core.CreateLanguage(
		"gql",
		parent,
		nil,
		gqlScripts,
		&gqlTpl,
		gqlTypes,
		"[%s]",
	)
}
