package _go_test

import (
	"github.com/jt05610/scaf/codegen"
	"github.com/jt05610/scaf/core"
	_go "github.com/jt05610/scaf/go"
	"github.com/jt05610/scaf/testData"
	"testing"
)

func TestGenerate(t *testing.T) {
	testData.RunTest(t, "testing", _go.Lang("testing"), func(system *core.System) error {
		gen := codegen.New("testing", _go.Lang("testing"))
		return gen.VisitSystem(system)
	})
}
