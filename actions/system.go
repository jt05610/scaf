package actions

import "github.com/jt0610/scaf/system"

func (h *Handler) CreateSystem(system *system.System) error {
	return h.SystemService.Create(system)
}
