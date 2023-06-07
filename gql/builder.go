package gql

import (
	"github.com/jt05610/scaf/builder"
	"github.com/jt05610/scaf/core"
	"os/exec"
)

var ParDir = builder.NewDir(func(m *core.Module) string { return "." }, []*builder.File{
	{
		Name:     func(m *core.Module) string { return "tools.go" },
		Template: toolsTemplate,
	},
	{
		Name:     func(m *core.Module) string { return "gqlgen.yml" },
		Template: gqlGenYmlTemplate,
	},
	{
		Name:     func(m *core.Module) string { return "server.go" },
		Template: serverTemplateInstance,
	},
})
var Dirs = builder.NewDir(func(*core.Module) string { return "graph" }, []*builder.File{
	{
		Name:     func(m *core.Module) string { return "schema.graphql" },
		Template: gqlTemplate,
	},
	{
		Name:     func(m *core.Module) string { return "schema.resolvers.go" },
		Template: schemaResolverTemplate,
	},
	{
		Name:     func(m *core.Module) string { return "resolver.go" },
		Template: resolverTemplate,
	},
})

func runGqlGen(m *core.Module) *exec.Cmd {
	cmd := exec.Command("go", "run", "github.com/99designs/gqlgen", "generate")
	cmd.Dir = m.Name
	return cmd
}

func NewBuilder() core.Visitor {
	ParDir.AddChild(Dirs)

	return builder.NewBuilder(
		builder.NewDirBuilder(ParDir),
		builder.NewRunner(builder.GoModTidy),
		builder.NewRunner(runGqlGen),
		builder.NewRunner(builder.GoFmt),
	)
}
