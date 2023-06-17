package builder_test

import (
	goContext "context"
	"github.com/jt05610/scaf/builder"
	"github.com/jt05610/scaf/cmd"
	"github.com/jt05610/scaf/context"
	"github.com/jt05610/scaf/core"
	"github.com/jt05610/scaf/lang"
	"github.com/jt05610/scaf/testData"
	"github.com/jt05610/scaf/zap"
	"testing"
)

func TestSupervisor_Start(t *testing.T) {
	t.Run("should start all modules without errors", func(t *testing.T) {
		parent := "testParent"
		mainCtx, cancel := goContext.WithCancel(context.Background())
		defer cancel()
		logger := zap.NewDev(mainCtx, parent, "supervisor")
		ctx := context.NewContext(logger)
		backend := lang.Go(parent)
		testData.RunTest(t, parent, backend, func(system *core.System) error {
			cmd.Gen(ctx, parent, system)
			supervisor := builder.NewSupervisor(parent, system)
			err := supervisor.Start(ctx)
			if err != nil {
				t.Errorf("failed to start modules: %v", err)
			}
			return err
		})
	})
}
