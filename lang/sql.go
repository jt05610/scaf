package lang

import (
	"embed"
	"github.com/jt05610/scaf/core"
)

//go:embed template/sql
var sqlTpl embed.FS

var sqlTypes = TypeMap{
	core.Int:    "integer",
	core.Float:  "double precision",
	core.String: "text",
	core.Bool:   "boolean",
	core.ID:     "uuid",
}

func SQL(parent string) *Language {
	return CreateLanguage(
		"sql",
		parent,
		nil,
		&Scripts{
			Map: map[string]string{
				"gen": "sqlc generate",
			},
		},
		&sqlTpl,
		sqlTypes,
		"%s[]",
	)
}
