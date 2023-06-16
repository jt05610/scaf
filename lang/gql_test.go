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

func TestGQL(t *testing.T) {
	gql := lang.GraphQL("testing")
	l := zap.NewDev(context.Background(), "testing", gql.Name)
	ctx := context.NewContext(l)
	testData.RunTest(t, "testing", gql, func(system *core.System) error {
		gen := codegen.New("testing", gql)
		return gen.VisitSystem(ctx, system)
	})

}
