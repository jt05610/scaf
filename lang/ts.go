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
	Init: `
npm install --save-dev --save-exact prettier
echo {}> .prettierrc.json
npm install --save-dev typescript @types/node
npx prettier --write .
`,
}

func TypeScript(parent string) *core.Language {
	return core.CreateLanguage(
		"ts",
		parent,
		tsScripts,
		&tsTpl,
		tsTypes,
		"%s[]",
	)
}
