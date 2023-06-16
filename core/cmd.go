package core

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
)

type Cmd struct {
	init  []func(m *Module) *exec.Cmd
	gen   []func(m *Module) *exec.Cmd
	start []func(m *Module) *exec.Cmd
	stop  []func(m *Module) *exec.Cmd
}

func (c *Cmd) Gen() []func(m *Module) *exec.Cmd {
	return c.gen
}

func (c *Cmd) Init() []func(m *Module) *exec.Cmd {
	return c.init
}

func (c *Cmd) Start() []func(m *Module) *exec.Cmd {
	return c.start
}

func (c *Cmd) Stop() []func(m *Module) *exec.Cmd {
	return c.stop
}

func CmdFuncs(parent, cc string) []func(m *Module) *exec.Cmd {
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
				return cmd
			})
		}
	}
	return cmds
}

func NewCmd(parent string, scripts *Scripts) *Cmd {
	return &Cmd{
		init:  CmdFuncs(parent, scripts.Init),
		gen:   CmdFuncs(parent, scripts.Gen),
		start: CmdFuncs(parent, scripts.Start),
		stop:  CmdFuncs(parent, scripts.Stop),
	}
}
