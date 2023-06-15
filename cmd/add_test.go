package cmd_test

import (
	"github.com/jt05610/scaf/builder"
	"github.com/jt05610/scaf/codegen"
	_go "github.com/jt05610/scaf/go"
	"github.com/jt05610/scaf/gql"
	"github.com/jt05610/scaf/proto"
	"github.com/jt05610/scaf/testData"
	"os"
	"testing"
)

func TestAdd(t *testing.T) {
	parent := "testing"
	err := os.RemoveAll(parent)
	if err != nil && !os.IsNotExist(err) {
		t.Fatal(err)
	}
	goLang := _go.Lang(parent)
	s := testData.HouseworkSystem(parent, goLang)
	GQL := gql.Lang(parent)
	protoBuf := proto.Lang(parent)
	bld := builder.NewBuilder(
		codegen.New("testing", gql.Lang("testing")),
		codegen.New("testing", proto.Lang("testing")),
		codegen.New("testing", goLang),
		builder.NewRunner(parent, protoBuf.Init()),
		builder.NewRunner(parent, goLang.Init()),
		builder.NewRunner(parent, GQL.Init()),
	)
	err = bld.VisitSystem(s)
	if err != nil {
		t.Error(err)
	}
}
