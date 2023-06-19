package core

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
)

type Cmd[T Storable] struct {
	init  []func(m T) *exec.Cmd
	gen   []func(m T) *exec.Cmd
	start []func(m T) *exec.Cmd
	stop  []func(m T) *exec.Cmd
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

func ModFuncs(parent, cc string, sub ...string) []func(m *Module) *exec.Cmd {
	lines := strings.Split(cc, "\n")
	cmds := make([]func(m *Module) *exec.Cmd, 0)
	for _, c := range lines {
		if c == "" {
			continue
		}
		t := template.Must(template.New(c).Parse(c))

		if len(strings.Split(c, " ")) > 0 {
			cmds = append(cmds, func(m *Module) *exec.Cmd {
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

func SysFuncs(parent, cc string, sub ...string) []func(s *System) *exec.Cmd {
	lines := strings.Split(cc, "\n")
	cmds := make([]func(s *System) *exec.Cmd, 0)
	for _, c := range lines {
		if c == "" {
			continue
		}
		t := template.Must(template.New(c).Parse(c))

		if len(strings.Split(c, " ")) > 0 {
			cmds = append(cmds, func(s *System) *exec.Cmd {
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
	Mod *Cmd[*Module]
	Sys *Cmd[*System]
}

func NewSysCmd(parent string, scripts *Scripts) (cmd *Cmd[*System]) {
	if scripts == nil {
		return nil
	}
	if scripts.WorkDir == "" {
		cmd = &Cmd[*System]{
			init:  SysFuncs(parent, scripts.Init),
			gen:   SysFuncs(parent, scripts.Gen),
			start: SysFuncs(parent, scripts.Start),
			stop:  SysFuncs(parent, scripts.Stop),
		}
	} else {
		cmd = &Cmd[*System]{
			init:  SysFuncs(parent, scripts.Init, scripts.WorkDir),
			gen:   SysFuncs(parent, scripts.Gen, scripts.WorkDir),
			start: SysFuncs(parent, scripts.Start, scripts.WorkDir),
			stop:  SysFuncs(parent, scripts.Stop, scripts.WorkDir),
		}
	}
	return
}

func NewModCmd(parent string, scripts *Scripts) (cmd *Cmd[*Module]) {
	if scripts == nil {
		return nil
	}
	if scripts.WorkDir == "" {
		cmd = &Cmd[*Module]{
			init:  ModFuncs(parent, scripts.Init),
			gen:   ModFuncs(parent, scripts.Gen),
			start: ModFuncs(parent, scripts.Start),
			stop:  ModFuncs(parent, scripts.Stop),
		}
	} else {
		cmd = &Cmd[*Module]{
			init:  ModFuncs(parent, scripts.Init, scripts.WorkDir),
			gen:   ModFuncs(parent, scripts.Gen, scripts.WorkDir),
			start: ModFuncs(parent, scripts.Start, scripts.WorkDir),
			stop:  ModFuncs(parent, scripts.Stop, scripts.WorkDir),
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
