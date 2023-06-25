package lang

import (
	"embed"
	"github.com/jt05610/scaf/core"
)

//go:embed template/proto
var protobufTpl embed.FS

var protobufTypes = TypeMap{
	core.Int:    "int32",
	core.Float:  "float",
	core.String: "string",
	core.ID:     "string",
	core.Bool:   "bool",
}

var protobufScripts = &Scripts{
	Map: map[string]string{
		"init": `

`,
		"gen": `
protoc --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. v{{.Version}}/{{.Name}}.proto
`,
	},
}

func Protobuf(parent string) *Language {
	return CreateLanguage(
		"proto",
		parent,
		nil,
		protobufScripts,
		&protobufTpl,
		protobufTypes,
		"repeated %s",
	)
}
