package core

type Visitor interface {
	Visit(m *Module) error
}
