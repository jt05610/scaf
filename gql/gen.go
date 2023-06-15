package gql

import (
	"embed"
	"github.com/jt05610/scaf/builder"
	"github.com/jt05610/scaf/core"
)

//go:embed template
var templates embed.FS

var initCmd = `
go mod tidy
go generate ./...
go fmt ./...
rm ./models_gen.go
`

var startCmd = `
{{.Name}} serve --port {{.Port}} 
`

var stopCmd = `
kill $(lsof -t -i:{{.Port}})
`

type lang struct {
	*builder.Cmd
}

func (l *lang) FS() *embed.FS {
	return &templates
}

var typeMap = map[core.BaseType]string{
	core.IntType:    "Int",
	core.FloatType:  "Float",
	core.StringType: "String",
	core.BoolType:   "Boolean",
}

func (l *lang) MapType(t core.BaseType) (s string, b bool) {
	s, b = typeMap[t]
	return s, b
}

func (l *lang) MakeArray(s string) string {
	return "[" + s + "]"
}

func Lang(parent string) core.Language {
	return &lang{
		Cmd: builder.NewCmd(parent, initCmd, startCmd, stopCmd),
	}
}
