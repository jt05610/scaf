package testData

import (
	"github.com/jt05610/scaf/core"
	"testing"
)

func HouseworkSystem(name string, backend *core.Language) *core.System {
	choreType := &core.Type{
		Name: "Chore",
		Fields: []*core.Field{
			{Name: "complete", Type: core.Bool},
			{Name: "description", Type: core.String},
		},
		Query: true,
	}
	mod := &core.Module{
		MetaData: &core.MetaData{
			Name:   "housework",
			Date:   "07 Jun 2023",
			Author: "Jonathan Taylor",
		},
		API: map[string]*core.API{
			"v1": {
				Name:     "housework",
				Version:  1,
				Date:     "07 Jun 2023",
				Author:   "Jonathan Taylor",
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
			},
		},
	}
	s := core.NewSystem(name, "Jonathan Taylor", "15 Jun 2023")
	err := s.AddModule(mod)
	if err != nil {
		panic(err)
	}
	return s
}

func RunTest(t *testing.T, parent string, backend *core.Language, f func(system *core.System) error) {
	s := HouseworkSystem(parent, backend)
	err := f(s)
	if err != nil {
		t.Error(err)
	}
}
