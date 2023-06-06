package builder

import (
	"github.com/jt05610/scaf/core"
	"os/exec"
)

func GoModInit(m *core.Module) *exec.Cmd {
	cmd := exec.Command("go", "mod", "init", m.Name)
	cmd.Dir = m.Name
	return cmd
}

func GoModTidy(m *core.Module) *exec.Cmd {
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = m.Name
	return cmd
}

func GoFmt(m *core.Module) *exec.Cmd {
	cmd := exec.Command("go", "fmt", "./...")
	cmd.Dir = m.Name
	return cmd
}
