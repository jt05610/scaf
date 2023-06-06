package core_test

import (
	"github.com/jt05610/scaf/core"
	"testing"
)

// utility function for creating a simple module
func createModule(name string, deps ...*core.Module) *core.Module {
	return &core.Module{Name: name, Deps: deps}
}

func TestIsAcyclic(t *testing.T) {
	modA := createModule("A")
	modB := createModule("B", modA)
	modC := createModule("C", modB)

	s := &core.System{
		Modules: []*core.Module{
			modA,
			modB,
			modC,
		},
	}

	c := core.NewChecker()

	if !c.IsAcyclic(s) {
		t.Errorf("Expected system to be acyclic")
	}

	// create a cycle
	modA.Deps = append(modA.Deps, modC)

	if c.IsAcyclic(s) {
		t.Errorf("Expected system to be cyclic")
	}
}
