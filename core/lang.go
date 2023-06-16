package core

import (
	"embed"
	"os/exec"
)

type Language interface {
	FS() *embed.FS
	MapType(t BaseType) (string, bool)
	MakeArray(s string) string
	Gen() []func(*Module) *exec.Cmd
	Start() []func(*Module) *exec.Cmd
	Stop() []func(*Module) *exec.Cmd
}
