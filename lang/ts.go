package lang

import (
	"embed"
	"github.com/jt05610/scaf/core"
)

//go:embed template/ts
var tsTpl embed.FS

var tsTypes = TypeMap{
	core.Int:    "number",
	core.Float:  "number",
	core.String: "string",
	core.ID:     "string",
	core.Bool:   "boolean",
}

var tsScripts = &Scripts{
	WorkDir: "ui",
}

func TypeScript(parent string) *Language {
	tsSysScripts := &Scripts{
		Map: map[string]string{
			"init": `
`,
			"gen": `
pnpm install
npx prettier --write .
`,
			"build": `
pnpm run build
`,
		},
	}
	return CreateLanguage(
		"ts",
		parent,
		tsSysScripts,
		tsScripts,
		&tsTpl,
		tsTypes,
		"%s[]",
	)
}
