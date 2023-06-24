package lang

import (
	"embed"
	"github.com/jt05610/scaf/core"
)

//go:embed template/gql
var gqlTpl embed.FS

var gqlTypes = TypeMap{
	core.Int:    "Int",
	core.Float:  "Float",
	core.String: "String",
	core.Bool:   "Boolean",
	core.ID:     "ID!",
}

func GraphQL(parent string) *Language {
	return CreateLanguage(
		"gql",
		parent,
		nil,
		nil,
		&gqlTpl,
		gqlTypes,
		"[%s]",
	)
}
