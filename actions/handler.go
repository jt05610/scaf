package actions

import (
	"github.com/jt0610/scaf/caddy"
	"github.com/jt0610/scaf/context"
	"github.com/jt0610/scaf/service"
	"github.com/jt0610/scaf/system"
	"github.com/jt0610/scaf/yaml"
)

// Handler is a type that combines context.Context with Services related to system.System and system.Module.
// The services are responsible for handling operations related to instances of system.System and system.Module.
type Handler struct {
	// Ctx is the context in which the handler operates. This can be used to handle cancellations, timeouts, and passing request-scoped values.
	Ctx context.Context
	// SystemService is a service.Service specialized to work with system.System.
	SystemService service.Service[*system.System]
	// CaddyService is a service.Service specialized to work with caddy.Caddyfile.
	CaddyService service.Service[*caddy.Caddyfile]
	// ModService is a service.Service specialized to work with system.Module.
	ModService service.Service[*system.Module]
}

// YamlHandler returns a new instance of Handler with its context set to the provided context and its services initialized to instances of yaml.YAMLService for system.System and system.Module.
// This can be used to handle operations related to system.System and system.Module in a context where YAML is used as the data interchange format.
func YamlHandler(ctx context.Context) *Handler {
	return &Handler{
		Ctx:           ctx,
		SystemService: yaml.NewYAMLService[*system.System](),
		CaddyService:  yaml.NewYAMLService[*caddy.Caddyfile](),
		ModService:    yaml.NewYAMLService[*system.Module](),
	}
}
