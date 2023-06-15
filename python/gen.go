package python

import (
	"embed"
	"github.com/jt05610/scaf/builder"
	"github.com/jt05610/scaf/core"
)

//go:embed template
var templates embed.FS

var initCmd = `
mkcert {{.Name}}.local localhost 127.0.0.1 ::1
mv {{.Name}}.local+3.pem {{.Name}}.local+3-key.pem ./cmd/.secrets
go mod tidy
go fmt ./...
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
	core.IntType:    "int",
	core.FloatType:  "float",
	core.StringType: "str",
	core.BoolType:   "bool",
}

func (l *lang) MapType(t core.BaseType) (s string, b bool) {
	s, b = typeMap[t]
	return s, b
}

func (l *lang) MakeArray(s string) string {
	return "List[" + s + "]"
}

func Lang(parent string) core.Language {
	return &lang{
		Cmd: builder.NewCmd(parent, initCmd, startCmd, stopCmd),
	}
}
