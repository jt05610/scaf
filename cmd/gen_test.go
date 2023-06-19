package cmd_test

import (
	"github.com/jt05610/scaf/cmd"
	"github.com/jt05610/scaf/context"
	"github.com/jt05610/scaf/testData"
	"github.com/jt05610/scaf/zap"
	"os"
	"testing"
)

func TestGen(t *testing.T) {
	parent := "testing"
	err := os.RemoveAll(parent)
	if err != nil && !os.IsNotExist(err) {
		t.Fatal(err)
	}
	s := testData.SCAFSystem(parent)
	logger := zap.NewDev(context.Background(), "testing", "gen_test")
	ctx := context.NewContext(logger)
	cmd.Gen(ctx, parent, s)
}
