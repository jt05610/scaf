package lang

import (
	"embed"
	"github.com/jt05610/scaf/core"
)

//go:embed template/ts
var tsTpl embed.FS

var tsTypes = &core.TypeMap{
	Int:    "number",
	Float:  "number",
	String: "string",
	Bool:   "boolean",
}

var tsScripts = &core.Scripts{
	WorkDir: "ui",
}

func TypeScript(parent string) *core.Language {
	tsSysScripts := &core.Scripts{
		Init: `
`,
		Gen: `
pnpm install
npx prettier --write .
`,
		Build: `
pnpm run build
`,
	}
	return core.CreateLanguage(
		"ts",
		parent,
		tsSysScripts,
		tsScripts,
		&tsTpl,
		tsTypes,
		"%s[]",
	)
}
