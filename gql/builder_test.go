package gql

import (
	"github.com/jt05610/scaf/core"
	"testing"
)

func TestBuilder(t *testing.T) {
	tests := []struct {
		name   string
		module *core.Module
	}{
		{
			name: "testMod",
			module: &core.Module{
				Name: "simple",
				Port: 8081,
				Date: "05 Jun 2023",
				Types: []*core.Type{
					{
						Name: "Example",
						Fields: []*core.Field{
							{Name: "exampleField", Type: "String"},
						},
					},
				},
				Funcs: []*core.Func{
					{
						Name: "exampleFunc",
						Params: []*core.Field{
							{Name: "exampleParam", Type: "String"},
						},
						Return: []*core.Field{
							{Name: "exampleReturn", Type: "String"},
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			b := NewBuilder()
			v := b.Visit(test.module)
			if v == nil {
			}
		})
	}
}
