package builder

import (
	"bytes"
	"fmt"
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

func Run(ctx context.Context, m *core.Module, cfs []func(m *core.Module) *exec.Cmd) error {
	for _, cf := range cfs {
		cmd := cf(m)
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr

		ctx.Logger.Info(cmd.String())
		err := cmd.Run()
		ctx.Logger.Info(cmd.String(), zap.String("output", out.String()))
		fmt.Println(stderr.String())
		if err != nil {
			ctx.Logger.Error(cmd.String(), zap.String("stderr", stderr.String()), zap.Error(err))
			fmt.Println(stderr.String())
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
		m.Version = api.Version
		err := Run(ctx, m, r.cfs)
		if err != nil {
			return err
		}
	}
	return nil
}

func NewRunner(parent string, cfs []func(m *core.Module) *exec.Cmd) core.Visitor {
	return &runner{cfs: cfs, parent: parent, seen: make(map[string]bool)}
}
