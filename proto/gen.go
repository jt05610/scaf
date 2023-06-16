package proto

import (
	"embed"
	"github.com/jt05610/scaf/builder"
	"github.com/jt05610/scaf/core"
)

//go:embed template
var templates embed.FS

var initCmd = `
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
protoc --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. v{{.Version}}/{{.Name}}.proto
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
	core.IntType:    "int32",
	core.FloatType:  "float",
	core.StringType: "string",
	core.BoolType:   "bool",
}

func (l *lang) MapType(t core.BaseType) (s string, b bool) {
	s, b = typeMap[t]
	return s, b
}

func (l *lang) MakeArray(s string) string {
	return "repeated " + s
}

func Lang(parent string) core.Language {
	return &lang{
		Cmd: builder.NewCmd(parent, initCmd, startCmd, stopCmd),
	}
}
