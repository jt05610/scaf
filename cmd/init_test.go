package cmd_test

import (
	"github.com/jt05610/scaf/cmd"
	"github.com/jt05610/scaf/context"
	"github.com/jt05610/scaf/zap"
	"os"
	"testing"
)

func TestInit(t *testing.T) {
	parent := "initTest"
	err := os.RemoveAll(parent)
	if err != nil && !os.IsNotExist(err) {
		t.Fatal(err)
	}
	logger := zap.NewDev(context.Background(), "testing", "gen_test")
	ctx := context.NewContext(logger)
	cmd.Init(ctx, parent)
}
