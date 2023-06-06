package builder

import (
	"github.com/jt05610/scaf/core"
	"os/exec"
)

type runner struct {
	cf   func(m *core.Module) *exec.Cmd
	seen map[string]bool
}

func (r *runner) Visit(m *core.Module) core.Visitor {
	if _, seen := r.seen[m.Name]; seen {
		return nil
	}
	r.seen[m.Name] = true
	cmd := r.cf(m)
	cmd.Dir = m.Name
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
	return r
}

func NewRunner(cf func(m *core.Module) *exec.Cmd) core.Visitor {
	return &runner{cf: cf, seen: make(map[string]bool)}
}
