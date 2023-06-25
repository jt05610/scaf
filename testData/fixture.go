package testData

import (
	"github.com/jt05610/scaf/core"
)

// Users creates a module for managing users
func Users() *core.Module {
	model := &core.Model{
		Name:   "User",
		Plural: "Users",
		Query:  true,
		Create: true,
		Update: true,
		Delete: true,
		Fields: []*core.Field{
			{
				Name:     "Name",
				Type:     core.String,
				Required: true,
				Create:   true,
			},
			{
				Name:     "Email",
				Type:     core.String,
				Required: true,
				Create:   true,
			},
		},
	}
	mod := core.NewModule("users", "Jonathan Taylor", "23 Jun 2023")
	modAPI := core.NewAPI("users", "Jonathan Taylor", "23 Jun 2023", core.Go)
	modAPI.AddModel(model)
	mod.AddAPI(modAPI)
	return mod
}

// APIs creates a module for managing APIs
func APIs() *core.Module {
	user := core.External("User", "User", "Users")

	dep := &core.Model{
		Name:   "Dependency",
		Plural: "Dependencies",
		Fields: []*core.Field{
			{
				Name:     "Name",
				Type:     core.String,
				Required: true,
				Create:   true,
			},
			{
				Name:     "Version",
				Type:     core.String,
				Required: true,
				Create:   true,
				Update:   true,
			},
		},
	}
	model := &core.Model{
		Name:   "Model",
		Plural: "Models",
		Fields: []*core.Field{
			{
				Name:     "Name",
				Type:     core.String,
				Required: true,
				Create:   true,
			},
			{
				Name:     "Description",
				Type:     core.String,
				Required: true,
				Create:   true,
			},
			{
				Name:     "Plural",
				Type:     core.String,
				Required: true,
				Create:   true,
			},
		},
	}
	field := &core.Model{
		Name:   "Field",
		Plural: "Fields",
		Fields: []*core.Field{
			{
				Name:     "Name",
				Type:     core.String,
				Required: true,
				Create:   true,
			},
			{
				Name:     "Description",
				Type:     core.String,
				Required: true,
				Create:   true,
			},
			{
				Name:     "Type",
				Type:     core.String,
				Required: true,
				Create:   true,
			},
			{
				Name:     "Required",
				Type:     core.Bool,
				Required: true,
				Create:   true,
			},
			{
				Name:   "IsArray",
				Type:   core.Bool,
				Update: true,
				Create: true,
			},
			{
				Name:   "Query",
				Type:   core.Bool,
				Create: true,
				Query:  true,
			},
			{
				Name:   "Create",
				Type:   core.Bool,
				Create: true,
				Query:  true,
			},
			{
				Name:   "Update",
				Type:   core.Bool,
				Create: true,
				Query:  true,
			},
			{
				Name:   "Delete",
				Type:   core.Bool,
				Create: true,
				Query:  true,
			},
			{
				Name:   "Subscribe",
				Type:   core.Bool,
				Create: true,
				Query:  true,
			},
		},
	}

	model.Fields = append(model.Fields, &core.Field{
		Name:    "Fields",
		Type:    field,
		IsArray: true,
	})

	function := &core.Model{
		Name:   "Function",
		Plural: "Functions",
		Fields: []*core.Field{
			{
				Name:     "Name",
				Type:     core.String,
				Required: true,
				Create:   true,
			},
			{
				Name:     "Description",
				Type:     core.String,
				Required: true,
				Create:   true,
				Update:   true,
			},
			{
				Name:    "Params",
				Type:    field,
				IsArray: true,
				Create:  true,
			},
			{
				Name:    "Returns",
				Type:    field,
				IsArray: true,
				Create:  true,
			},
			{
				Name:   "UserCode",
				Type:   core.String,
				Create: true,
				Update: true,
			},
		},
	}
	api := &core.Model{
		Name:   "Api",
		Plural: "Apis",
		Query:  true,
		Create: true,
		Fields: []*core.Field{
			{
				Name:     "Name",
				Type:     core.String,
				Required: true,
				Create:   true,
			},
			{
				Name: "Running",
				Type: core.Bool,
			},
			user,
			{
				Name: "Date",
				Type: core.String,
			},
			{
				Name: "Version",
				Type: core.String,
			},
			{
				Name:     "Description",
				Type:     core.String,
				Required: true,
				Create:   true,
				Update:   true,
			},
			{
				Name:    "Functions",
				Type:    function,
				IsArray: true,
				Create:  true,
			},
			{
				Name:    "Models",
				Type:    model,
				IsArray: true,
				Create:  true,
			},
			{
				Name:     "Language",
				Type:     core.String,
				Required: true,
				Create:   true,
			},
			{
				Name:     "Dependencies",
				Type:     dep,
				IsArray:  true,
				Required: true,
				Create:   true,
			},
		},
	}

	ret := core.NewModule("apis", "Jonathan Taylor", "23 Jun 2023")
	modAPI := core.NewAPI("apis", "Jonathan Taylor", "23 Jun 2023", core.Go)
	modAPI.AddModel(api)

	genFunc := &core.Func{
		Name:        "Generate",
		Description: "Creates a new module and protobuf APIs for any dependencies",
		Params: []*core.Field{
			{
				Name: "module",
				Type: core.ID,
			},
		},
		Return: []*core.Field{
			{
				Name: "message",
				Type: core.String,
			},
			{
				Name: "url",
				Type: core.String,
			},
		},
	}
	startFunc := &core.Func{
		Name:        "Start",
		Description: "Starts the module's servers",
		Params: []*core.Field{
			{
				Name: "module",
				Type: core.ID,
			},
		},
		Return: []*core.Field{
			{
				Name: "message",
				Type: core.String,
			},
		},
	}
	stopFunc := &core.Func{
		Name:        "Stop",
		Description: "Stops the module's servers",
		Params: []*core.Field{
			{
				Name: "module",
				Type: core.ID,
			},
		},
		Return: []*core.Field{
			{
				Name: "message",
				Type: core.String,
			},
		},
	}
	for _, _func := range []*core.Func{genFunc, startFunc, stopFunc} {
		modAPI.AddFunc(_func)
	}
	modAPI.Description = "The APIs module is used to create and manage APIs"
	ret.AddAPI(modAPI)
	return ret
}

// Modules creates a module for creating modules
func Modules() *core.Module {
	user := core.External("User", "User", "Users")
	api := core.External("APIs", "API", "APIs")
	api.IsArray = true
	api.Update = true
	api.Delete = true
	module := &core.Model{
		Name:   "Module",
		Plural: "Modules",
		Query:  true,
		Fields: []*core.Field{
			{
				Name:     "Name",
				Type:     core.String,
				Required: true,
				Create:   true,
				Query:    true,
				Update:   true,
			},
			{
				Name:     "Description",
				Type:     core.String,
				Create:   true,
				Required: true,
				Query:    true,
				Update:   true,
			},
			api,
			user,
		},
	}

	ret := core.NewModule("modules", "Jonathan Taylor", "23 Jun 2023")
	modAPI := core.NewAPI("modules", "Jonathan Taylor", "23 Jun 2023", core.Go)
	modAPI.AddModel(module)
	modAPI.Description = "The module for creating modules"
	ret.AddAPI(modAPI)
	return ret
}

func SCAFSystem(name string) *core.System {
	s := core.NewSystem(name, "The core system for scaf", "Jonathan Taylor", "18 Jun 2023")
	err := s.AddModule(APIs())
	if err != nil {
		panic(err)
	}
	err = s.AddModule(Users())
	if err != nil {
		panic(err)
	}
	err = s.AddModule(Modules())
	if err != nil {
		panic(err)
	}
	return s
}
