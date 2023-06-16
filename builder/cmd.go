package builder

import (
	"github.com/jt05610/scaf/core"
	"os/exec"
)

type Cmd struct {
	init  []func(m *core.Module) *exec.Cmd
	start []func(m *core.Module) *exec.Cmd
	stop  []func(m *core.Module) *exec.Cmd
}

func (c *Cmd) Gen() []func(m *core.Module) *exec.Cmd {
	return c.init
}

func (c *Cmd) Start() []func(m *core.Module) *exec.Cmd {
	return c.start
}

func (c *Cmd) Stop() []func(m *core.Module) *exec.Cmd {
	return c.stop
}

func NewCmd(parent, init, start, stop string) *Cmd {
	return &Cmd{
		init:  CmdFuncs(parent, init),
		start: CmdFuncs(parent, start),
		stop:  CmdFuncs(parent, stop),
	}
}
