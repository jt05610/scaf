package builder

import (
	"bytes"
	"github.com/jt05610/scaf/context"
	"github.com/jt05610/scaf/core"
	"go.uber.org/zap"
	"os/exec"
)

type runner struct {
	cfs    []func(m *core.Module) *exec.Cmd
	parent string
	seen   map[string]bool
}

func (r *runner) VisitSystem(ctx context.Context, s *core.System) error {
	ctx.Logger.Debug("Running commands", zap.String("system", s.Name))
	r.seen = make(map[string]bool)
	for _, m := range s.Modules {
		if err := r.VisitModule(ctx, m); err != nil {
			return err
		}
	}
	return nil
}

func (r *runner) VisitModule(ctx context.Context, m *core.Module) error {
	if _, seen := r.seen[m.Name]; seen {
		return nil
	}
	r.seen[m.Name] = true
	for _, api := range m.API {
		for _, cf := range r.cfs {
			m.Version = api.Version
			cmd := cf(m)
			var out bytes.Buffer
			var stderr bytes.Buffer
			cmd.Stdout = &out
			cmd.Stderr = &stderr

			err := cmd.Run()
			ctx.Logger.Info(cmd.String(), zap.String("output", out.String()))
			if err != nil {
				ctx.Logger.Error(cmd.String(), zap.String("stderr", stderr.String()), zap.Error(err))
				return err
			}
		}
	}

	return nil
}
func NewRunner(parent string, cfs []func(m *core.Module) *exec.Cmd) core.Visitor {
	return &runner{cfs: cfs, parent: parent, seen: make(map[string]bool)}
}
