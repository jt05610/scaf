package builder

import (
	"bytes"
	"fmt"
	"github.com/jt05610/scaf/context"
	"github.com/jt05610/scaf/core"
	"github.com/jt05610/scaf/lang"
	"go.uber.org/zap"
	"os/exec"
	"sync"
)

type runner struct {
	cmd     *lang.CmdSet
	parent  string
	lang    core.Lang
	seen    map[string]bool
	mu      sync.Mutex
	scripts []ModScript
}

func (r *runner) VisitSystem(ctx context.Context, s *core.System) error {
	ctx.Logger.Debug("Running commands", zap.String("system", s.Name))
	r.seen = make(map[string]bool)
	var wg sync.WaitGroup
	for _, m := range s.Modules {
		wg.Add(1)
		go func(m *core.Module) {
			defer wg.Done()
			if err := r.VisitModule(ctx, m); err != nil {
				panic(err)
			}
		}(m)
	}
	wg.Wait()
	if r.cmd.Sys != nil {
		for _, script := range r.scripts {
			switch script {
			case ModScriptInit:
				ctx.Logger.Debug("Running system init commands", zap.String("system", s.Name))
				err := SysRun(ctx, s, r.cmd.Sys.Init())
				if err != nil {
					return err
				}
			case ModScriptStop:
				ctx.Logger.Debug("Running system exit commands", zap.String("system", s.Name))
				err := SysRun(ctx, s, r.cmd.Sys.Stop())
				if err != nil {
					return err
				}
			case ModScriptGen:
				ctx.Logger.Debug("Running system gen commands", zap.String("system", s.Name))
				err := SysRun(ctx, s, r.cmd.Sys.Gen())
				if err != nil {
					return err
				}
			case modScriptStart:
				ctx.Logger.Debug("Running system start commands", zap.String("system", s.Name))
				err := SysRun(ctx, s, r.cmd.Sys.Start())
				if err != nil {
					return err
				}
			case ModScriptBuild:
				ctx.Logger.Debug("Running system build commands", zap.String("system", s.Name))
				err := SysRun(ctx, s, r.cmd.Sys.Build())
				if err != nil {
					return err
				}
			}
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

func SysRun(ctx context.Context, s *core.System, cfs []func(m *core.System) *exec.Cmd) error {
	for _, cf := range cfs {
		cmd := cf(s)
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
	r.mu.Lock()
	if _, seen := r.seen[m.Name]; seen {
		return nil
	}
	r.seen[m.Name] = true
	r.mu.Unlock()
	ctx.Logger.Debug("Running commands", zap.String("module", m.Name))
	for _, api := range m.APIs() {
		if len(api.Funcs) == 0 && r.lang == core.Protobuf {
			continue
		}
		m.Version = api.Version
		for _, s := range r.scripts {
			var err error
			switch s {
			case ModScriptGen:
				err = Run(ctx, m, r.cmd.Mod.Gen())
			case ModScriptInit:
				err = Run(ctx, m, r.cmd.Mod.Init())
			case modScriptStart:
				err = Run(ctx, m, r.cmd.Mod.Start())
			case ModScriptStop:
				err = Run(ctx, m, r.cmd.Mod.Stop())
			case ModScriptBuild:
				err = Run(ctx, m, r.cmd.Mod.Build())
			}
			if err != nil {
				return err
			}
		}
	}
	return nil
}

type ModScript int

const (
	ModScriptGen ModScript = iota
	ModScriptInit
	ModScriptBuild
	modScriptStart
	ModScriptStop
)

func NewRunner(parent string, language core.Lang, set *lang.CmdSet, scripts ...ModScript) core.Visitor {
	return &runner{
		cmd:     set,
		parent:  parent,
		lang:    language,
		seen:    make(map[string]bool),
		scripts: scripts,
	}
}
