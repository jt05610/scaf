package lang

import (
	"bytes"
	"github.com/jt05610/scaf/core"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
)

type Cmd[T core.Storable] struct {
	init  []func(m T) *exec.Cmd
	gen   []func(m T) *exec.Cmd
	start []func(m T) *exec.Cmd
	stop  []func(m T) *exec.Cmd
	build []func(m T) *exec.Cmd
}

func (c *Cmd[T]) Gen() []func(m T) *exec.Cmd {
	return c.gen
}

func (c *Cmd[T]) Init() []func(m T) *exec.Cmd {
	return c.init
}

func (c *Cmd[T]) Start() []func(m T) *exec.Cmd {
	return c.start
}

func (c *Cmd[T]) Stop() []func(m T) *exec.Cmd {
	return c.stop
}

func (c *Cmd[T]) Build() []func(m T) *exec.Cmd {
	return c.build
}

func ModFuncs(parent, cc string, sub ...string) []func(m *core.Module) *exec.Cmd {
	lines := strings.Split(cc, "\n")
	cmds := make([]func(m *core.Module) *exec.Cmd, 0)
	for _, c := range lines {
		if c == "" {
			continue
		}
		t := template.Must(template.New(c).Parse(c))

		if len(strings.Split(c, " ")) > 0 {
			cmds = append(cmds, func(m *core.Module) *exec.Cmd {
				var buf bytes.Buffer
				err := t.Execute(&buf, m)
				if err != nil {
					panic(err)
				}
				c = buf.String()
				vars := os.Environ()
				pth := filepath.Join(os.Getenv("GOPATH"), "bin")
				vars = append(vars, "PATH="+os.Getenv("PATH")+":"+pth)
				args := strings.Split(c, " ")
				cmd := exec.Command(args[0], args[1:]...)
				cmd.Env = append(cmd.Env, vars...)
				cmd.Dir = filepath.Join(parent, m.Name)
				if len(sub) > 0 {
					cmd.Dir = filepath.Join(cmd.Dir, sub[0])
				}
				return cmd
			})
		}
	}
	return cmds
}

func SysFuncs(parent, cc string, sub ...string) []func(s *core.System) *exec.Cmd {
	lines := strings.Split(cc, "\n")
	cmds := make([]func(s *core.System) *exec.Cmd, 0)
	for _, c := range lines {
		if c == "" {
			continue
		}
		t := template.Must(template.New(c).Parse(c))

		if len(strings.Split(c, " ")) > 0 {
			cmds = append(cmds, func(s *core.System) *exec.Cmd {
				var buf bytes.Buffer
				err := t.Execute(&buf, s)
				if err != nil {
					panic(err)
				}
				c = buf.String()
				vars := os.Environ()
				pth := filepath.Join(os.Getenv("GOPATH"), "bin")
				vars = append(vars, "PATH="+os.Getenv("PATH")+":"+pth)
				args := strings.Split(c, " ")
				cmd := exec.Command(args[0], args[1:]...)
				cmd.Env = append(cmd.Env, vars...)
				if len(sub) > 0 {
					cmd.Dir = filepath.Join(parent, sub[0])
				} else {
					cmd.Dir = parent
				}
				return cmd
			})
		}
	}
	return cmds
}

type CmdSet struct {
	Mod *Cmd[*core.Module]
	Sys *Cmd[*core.System]
}

func NewSysCmd(parent string, scripts *Scripts) (cmd *Cmd[*core.System]) {
	if scripts == nil {
		return nil
	}
	if scripts.WorkDir == "" {
		cmd = &Cmd[*core.System]{
			init:  SysFuncs(parent, scripts.Map["init"]),
			gen:   SysFuncs(parent, scripts.Map["gen"]),
			start: SysFuncs(parent, scripts.Map["start"]),
			stop:  SysFuncs(parent, scripts.Map["stop"]),
			build: SysFuncs(parent, scripts.Map["build"]),
		}
	} else {
		cmd = &Cmd[*core.System]{
			init:  SysFuncs(parent, scripts.Map["init"], scripts.WorkDir),
			gen:   SysFuncs(parent, scripts.Map["gen"], scripts.WorkDir),
			start: SysFuncs(parent, scripts.Map["start"], scripts.WorkDir),
			stop:  SysFuncs(parent, scripts.Map["stop"], scripts.WorkDir),
			build: SysFuncs(parent, scripts.Map["build"], scripts.WorkDir),
		}
	}
	return
}

func NewModCmd(parent string, scripts *Scripts) (cmd *Cmd[*core.Module]) {
	if scripts == nil {
		return nil
	}
	if scripts.WorkDir == "" {
		cmd = &Cmd[*core.Module]{
			init:  ModFuncs(parent, scripts.Map["init"]),
			gen:   ModFuncs(parent, scripts.Map["gen"]),
			start: ModFuncs(parent, scripts.Map["start"]),
			stop:  ModFuncs(parent, scripts.Map["stop"]),
			build: ModFuncs(parent, scripts.Map["build"]),
		}
	} else {
		cmd = &Cmd[*core.Module]{
			init:  ModFuncs(parent, scripts.Map["init"], scripts.WorkDir),
			gen:   ModFuncs(parent, scripts.Map["gen"], scripts.WorkDir),
			start: ModFuncs(parent, scripts.Map["start"], scripts.WorkDir),
			stop:  ModFuncs(parent, scripts.Map["stop"], scripts.WorkDir),
			build: ModFuncs(parent, scripts.Map["build"], scripts.WorkDir),
		}
	}
	return
}

func NewCmdSet(parent string, sysScripts, modScripts *Scripts) *CmdSet {
	return &CmdSet{
		Mod: NewModCmd(parent, modScripts),
		Sys: NewSysCmd(parent, sysScripts),
	}
}
