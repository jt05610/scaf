package protobuf

import (
	"github.com/jt05610/scaf/core"
	"os"
	"testing"
)

func TestBuilder(t *testing.T) {
	tests := []struct {
		name   string
		module *core.Module
	}{
		{
			name: "housework",
			module: &core.Module{
				Name: "housework",
				Port: 8081,
				Date: "06 Jun 2023",
				Types: []*core.Type{
					{
						Name: "Chore",
						Fields: []*core.Field{
							{Name: "complete", Type: core.Bool},
							{Name: "description", Type: core.String},
						},
						HasPlural: true,
					},
				},
				Funcs: []*core.Func{
					{
						Name: "Add",
						Params: []*core.Field{
							{Name: "chores", Type: core.Custom, CustomType: "Chores"},
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
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_ = os.RemoveAll(test.module.Name)
			b := NewBuilder()
			v := b.Visit(test.module)
			if v == nil {
			}
		})
	}
}
