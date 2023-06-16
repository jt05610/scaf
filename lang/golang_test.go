package lang_test

import (
	"github.com/jt05610/scaf/codegen"
	"github.com/jt05610/scaf/context"
	"github.com/jt05610/scaf/core"
	. "github.com/jt05610/scaf/lang"
	"github.com/jt05610/scaf/testData"
	"github.com/jt05610/scaf/zap"
	"testing"
)

func TestGo(t *testing.T) {
	lang := Go("testing")
	l := zap.NewDev(context.Background(), "testing", lang.Name)
	ctx := context.NewContext(l)
	testData.RunTest(t, "testing", lang, func(system *core.System) error {
		gen := codegen.New("testing", lang)
		return gen.VisitSystem(ctx, system)
	})
}
