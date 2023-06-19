package testData

import (
	"github.com/jt05610/scaf/core"
	"github.com/jt05610/scaf/lang"
	"testing"
)

// Creator creates a module for that is used for creating modules
func Creator(parent string) *core.Module {
	t := &core.Type{
		Name:   "Type",
		Plural: "Types",
		Fields: []*core.Field{
			{Name: "Name", Type: core.String},
			{Name: "Description", Type: core.String},
			{Name: "Plural", Type: core.String},
			{Name: "Query", Type: core.Bool},
			{Name: "Mutate", Type: core.Bool},
			{Name: "Subscribe", Type: core.Bool},
		},
		Query:     false,
		Mutate:    false,
		Subscribe: false,
	}
	field := &core.Type{
		Name:   "Field",
		Plural: "Fields",
		Fields: []*core.Field{
			{Name: "Name", Type: core.String},
			{Name: "Description", Type: core.String},
			{Name: "Required", Type: core.Bool},
			{Name: "Type", Type: t},
		},
		Query:     false,
		Mutate:    false,
		Subscribe: false,
	}

	t.Fields = append(t.Fields, &core.Field{Name: "Fields", Type: core.Array(field)})

	f := &core.Type{
		Name:   "Func",
		Plural: "Funcs",
		Fields: []*core.Field{
			{Name: "Name", Type: core.String},
			{Name: "Description", Type: core.String},
			{Name: "Params", Type: core.Array(field)},
			{Name: "Returns", Type: core.Array(field)},
		},
		Query:     false,
		Mutate:    false,
		Subscribe: false,
	}
	portMap := &core.Type{
		Name: "PortMap",
		Fields: []*core.Field{
			{Name: "UI", Type: core.Int},
			{Name: "GQL", Type: core.Int},
			{Name: "RPC", Type: core.Int},
		},
		Query:     false,
		Mutate:    false,
		Subscribe: false,
	}
	metaData := &core.Type{
		Name: "MetaData",
		Fields: []*core.Field{
			{Name: "ID", Type: core.String},
			{Name: "Rev", Type: core.String},
			{Name: "Name", Type: core.String},
			{Name: "Description", Type: core.String},
			{Name: "Author", Type: core.String},
			{Name: "Date", Type: core.String},
			{Name: "PortMap", Type: portMap},
		},
		Query:     false,
		Mutate:    false,
		Subscribe: false,
	}

	api := &core.Type{
		Name:   "API",
		Plural: "APIs",
		Fields: []*core.Field{
			{Name: "Name", Type: core.String},
			{Name: "PortMap", Type: portMap},
			{Name: "Author", Type: core.String},
			{Name: "Version", Type: core.String},
			{Name: "Language", Type: core.String},
			{Name: "Date", Type: core.String},
			{Name: "Types", Type: core.Array(t)},
			{Name: "Funcs", Type: core.Array(f)},
		},
		Query:     false,
		Mutate:    false,
		Subscribe: false,
	}
	module := &core.Type{
		Name:   "Module",
		Plural: "Modules",
		Fields: []*core.Field{
			{Name: "MetaData", Type: metaData},
			{Name: "Version", Type: core.Int},
			{Name: "APIs", Type: core.Array(api)},
		},
		Query:     true,
		Mutate:    true,
		Subscribe: false,
	}

	api.Fields = append(api.Fields, &core.Field{Name: "Deps", Type: core.Array(module)})
	types := []*core.Type{
		portMap,
		metaData,
		t,
		field,
		f,
		api,
		module,
	}

	return &core.Module{
		MetaData: &core.MetaData{
			Name:        "creator",
			Date:        "18 Jun 2023",
			Description: "scaf is a tool for creating systems",
			Author:      "Jonathan Taylor",
		},
		API: map[string]*core.API{
			"v1": {
				Name:     "creator",
				Version:  1,
				Date:     "18 Jun 2023",
				Author:   "Jonathan Taylor",
				Language: lang.Go(parent),
				Types:    types,
				Funcs:    []*core.Func{},
			},
		},
	}
}

func SCAFSystem(name string) *core.System {
	s := core.NewSystem(name, "The core system for scaf", "Jonathan Taylor", "18 Jun 2023")
	err := s.AddModule(Creator(name))
	if err != nil {
		panic(err)
	}
	return s
}

func RunTest(t *testing.T, parent string, f func(system *core.System) error) {
	s := SCAFSystem(parent)
	err := f(s)
	if err != nil {
		t.Error(err)
	}
}
