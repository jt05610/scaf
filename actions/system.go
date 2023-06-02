package actions

import (
	"github.com/jt0610/scaf/context"
	"github.com/jt0610/scaf/system"
	"github.com/jt0610/scaf/wizard"
	"io"
)

// CreateSystem is a method of Handler that facilitates the creation of a new system.
// It uses a wizard to interactively gather the necessary information from the user, populates the metadata of the system,
// and then uses the SystemService associated with the Handler to write the new system data to the provided io.Writer.
func (h *Handler) CreateSystem(ctx context.Context, writer io.Writer) (*system.System, error) {
	w := wizard.Wizard[system.System]{}
	sys, err := w.Run(ctx)
	if err != nil {
		return nil, err
	}
	err = sys.FillMetadata()
	if err != nil {
		return nil, err
	}
	return sys, h.SystemService.Flush(writer, sys)
}
