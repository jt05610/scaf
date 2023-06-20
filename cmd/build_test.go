package cmd_test

import (
	"github.com/jt05610/scaf/cmd"
	"github.com/jt05610/scaf/context"
	"github.com/jt05610/scaf/testData"
	"github.com/jt05610/scaf/zap"
	"testing"
)

func TestBuild(t *testing.T) {
	parent := "core"
	s := testData.SCAFSystem(parent)
	logger := zap.NewDev(context.Background(), "core", "build_test")
	ctx := context.NewContext(logger)
	cmd.Build(ctx, parent, s)
}
