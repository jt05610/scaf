package lang_test

import (
	"github.com/jt05610/scaf/codegen"
	"github.com/jt05610/scaf/context"
	"github.com/jt05610/scaf/core"
	"github.com/jt05610/scaf/lang"
	"github.com/jt05610/scaf/testData"
	"github.com/jt05610/scaf/zap"
	"testing"
)

func TestPython(t *testing.T) {
	py := lang.Python("testing")
	l := zap.NewDev(context.Background(), "testing", py.Name)
	ctx := context.NewContext(l)
	testData.RunTest(t, "testing", py, func(system *core.System) error {
		gen := codegen.New("testing", py)
		return gen.VisitSystem(ctx, system)
	})

}
