package builder

import (
	"bytes"
	"fmt"
	"github.com/jt05610/scaf/core"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
)

type runner struct {
	cfs    []func(m *core.Module) *exec.Cmd
	parent string
	seen   map[string]bool
}

func (r *runner) VisitSystem(s *core.System) error {
	r.seen = make(map[string]bool)
	for _, m := range s.Modules {
		if err := r.VisitModule(m); err != nil {
			return err
		}
	}
	return nil
}

func (r *runner) VisitModule(m *core.Module) error {
	if _, seen := r.seen[m.Name]; seen {
		return nil
	}
	r.seen[m.Name] = true
	for _, cf := range r.cfs {
		cmd := cf(m)
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr

		err := cmd.Run()
		fmt.Printf("%s\n", cmd.String())
		if err != nil {
			fmt.Println(stderr.String())
			return err
		} else {
			if out.String() != "" {
				fmt.Println(out.String())
			}
		}
	}
	return nil
}

func CmdFuncs(parent, cc string) []func(m *core.Module) *exec.Cmd {
	lines := strings.Split(cc, "\n")
	cmds := make([]func(m *core.Module) *exec.Cmd, 0)
	for _, c := range lines {
		if c == "" {
			continue
		}
		t := template.Must(template.New("cmd").Parse(c))

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
				return cmd
			})
		}
	}
	return cmds
}

func NewRunner(parent string, cfs []func(m *core.Module) *exec.Cmd) core.Visitor {
	return &runner{cfs: cfs, parent: parent, seen: make(map[string]bool)}
}
