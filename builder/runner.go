package builder

import (
	"bytes"
	"fmt"
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

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	fmt.Printf("%s\n", cmd.String())
	if err != nil {
		fmt.Printf("Error: %s\n", stderr.String())
		panic(err)
	} else {
		if out.String() != "" {
			fmt.Println(out.String())
		}
	}

	return r
}

func NewRunner(cf func(m *core.Module) *exec.Cmd) core.Visitor {
	return &runner{cf: cf, seen: make(map[string]bool)}
}
