package builder

import (
	"errors"
	"fmt"
	"github.com/jt05610/scaf/context"
	"github.com/jt05610/scaf/core"
	uz "go.uber.org/zap"
)

type Supervisor struct {
	*core.System
	parent string
	done   chan struct{}
	seen   map[string]bool
}

func NewSupervisor(parent string, s *core.System) *Supervisor {
	return &Supervisor{parent: parent, System: s}
}

func (s *Supervisor) Start(ctx context.Context) error {
	s.done = make(chan struct{})
	ctx.Logger.Info("Starting system", uz.String("system", s.Name))
	for _, m := range s.Modules {
		if _, seen := s.seen[m.Name]; seen {
			continue
		}
		s.seen[m.Name] = true
		if err := s.VisitModule(ctx, m); err != nil {
			return err
		}
	}
	return nil
}

func (s *Supervisor) VisitModule(ctx context.Context, m *core.Module) error {
	v := fmt.Sprintf("v%d", m.Version)
	api := m.API[v]
	if api == nil {
		ctx.Logger.Error("API not found for version", uz.String("version", v))
		return errors.New("API not found for version")
	}
	lang := api.Language
	if lang == nil {
		ctx.Logger.Error("Language not found for API", uz.String("api", api.Name))
		return errors.New("language not found for API")
	}
	ctx.Logger.Info("Starting module", uz.String("module", m.Name))
	go func() {
		go func() {
			err := Run(ctx, m, lang.Mod.Start())
			if err != nil {
				ctx.Logger.Error("Error running module", uz.String("module", m.Name), uz.Error(err))
			}
		}()
		<-s.done
	}()

	return nil
}
