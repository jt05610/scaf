package core

import "github.com/jt05610/scaf/context"

type Visitor interface {
	VisitSystem(ctx context.Context, s *System) error
	VisitModule(ctx context.Context, m *Module) error
}
