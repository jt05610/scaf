package actions

import (
	"github.com/jt0610/scaf/context"
	"github.com/jt0610/scaf/system"
	"github.com/jt0610/scaf/wizard"
	"io"
)

func (h *Handler) CreateSystem(ctx context.Context, writer io.Writer) error {
	w := wizard.Wizard[system.System]{}
	sys, err := w.Run(ctx)
	if err != nil {
		return err
	}
	return h.SystemService.Flush(writer, sys)
}
