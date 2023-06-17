package lang

import (
	"embed"
	"github.com/jt05610/scaf/core"
)

//go:embed template/proto
var protobufTpl embed.FS

var protobufTypes = &core.TypeMap{
	Int:    "int32",
	Float:  "float",
	String: "string",
	Bool:   "bool",
}

var protobufScripts = &core.Scripts{
	Init: `
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
`,
	Gen: `
protoc --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. v{{.Version}}/{{.Name}}.proto
`,
}

func Protobuf(parent string) *core.Language {
	return core.CreateLanguage(
		"proto",
		parent,
		protobufScripts,
		&protobufTpl,
		protobufTypes,
		"repeated %s",
	)
}
