package actions

import (
	"github.com/jt0610/scaf/context"
	"github.com/jt0610/scaf/service"
	"github.com/jt0610/scaf/system"
	"github.com/jt0610/scaf/yaml"
)

type Handler struct {
	Ctx           context.Context
	SystemService service.Service[system.System]
	ModService    service.Service[system.Module]
}

func YamlHandler(ctx context.Context) *Handler {
	return &Handler{
		Ctx:           ctx,
		SystemService: yaml.NewYAMLService[system.System](),
		ModService:    yaml.NewYAMLService[system.Module](),
	}
}
