package testData

import (
	"github.com/jt05610/scaf/core"
	"testing"
)

func HouseworkSystem(name string, backend core.Language) *core.System {
	choreType := &core.Type{
		Name: "Chore",
		Fields: []*core.Field{
			{Name: "complete", Type: core.Bool},
			{Name: "description", Type: core.String},
		},
		Query: true,
	}
	mod := &core.Module{
		Name:     "housework",
		Port:     8081,
		Date:     "07 Jun 2023",
		Author:   "Jonathan Taylor",
		Version:  1,
		Language: backend,
		Types: []*core.Type{
			choreType,
		},
		Funcs: []*core.Func{
			{
				Name: "Add",
				Params: []*core.Field{
					{Name: "chores", Type: core.Array(choreType)},
				},
				Return: []*core.Field{
					{Name: "message", Type: core.String},
				},
			},
			{
				Name: "Complete",
				Params: []*core.Field{
					{Name: "choreNumber", Type: core.Int},
				},
				Return: []*core.Field{
					{Name: "message", Type: core.String},
				},
			},
		},
	}
	return &core.System{
		Author:  "Jonathan Taylor",
		Date:    "15 Jun 2023",
		Name:    name,
		GQLPort: 8080,
		Modules: []*core.Module{mod},
	}
}

func RunTest(t *testing.T, parent string, backend core.Language, f func(system *core.System) error) {
	s := HouseworkSystem(parent, backend)
	err := f(s)
	if err != nil {
		t.Error(err)
	}
}
