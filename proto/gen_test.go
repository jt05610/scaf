package proto_test

import (
	"github.com/jt05610/scaf/codegen"
	"github.com/jt05610/scaf/core"
	"github.com/jt05610/scaf/proto"
	"github.com/jt05610/scaf/python"
	"github.com/jt05610/scaf/testData"
	"testing"
)

func TestGenerate(t *testing.T) {
	testData.RunTest(t, "testing", python.Lang("testing"), func(system *core.System) error {
		gen := codegen.New("testing", proto.Lang("testing"))
		return gen.VisitSystem(system)
	})
}
