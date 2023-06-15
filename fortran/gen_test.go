package fortran_test

import (
	"github.com/jt05610/scaf/codegen"
	"github.com/jt05610/scaf/core"
	"github.com/jt05610/scaf/fortran"
	"github.com/jt05610/scaf/testData"
	"github.com/jt05610/scaf/ts"
	"testing"
)

func TestGenerate(t *testing.T) {
	testData.RunTest(t, "testing", fortran.Lang("testing"), func(system *core.System) error {
		gen := codegen.New("testing", ts.Lang("testing"))
		return gen.VisitSystem(system)
	})
}
