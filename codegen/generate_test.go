package codegen

import (
	"github.com/jt05610/scaf/builder"
	"github.com/jt05610/scaf/core"
	"testing"
)

func TestGenerate(t *testing.T) {
	choreType := &core.Type{
		Name: "Chore",
		Fields: []*core.Field{
			{Name: "complete", Type: core.Bool},
			{Name: "description", Type: core.String},
		},
		Query: true,
	}
	mod := &core.Module{
		Name:    "housework",
		Port:    8081,
		Date:    "07 Jun 2023",
		Author:  "Jonathan Taylor",
		Version: 1,
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
	gen := builder.NewBuilder(Generator("testing"), builder.NewRunner("testing", Commands))
	err := gen.Visit(mod)
	if err != nil {
		t.Error(err)
	}
}
