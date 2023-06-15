package core

type Visitor interface {
	VisitSystem(s *System) error
	VisitModule(m *Module) error
}
