package builder

import "github.com/jt05610/scaf/core"

type Builder struct {
	Visitors []core.Visitor
	seen     map[string]bool
}

func (b *Builder) VisitSystem(s *core.System) error {
	for _, v := range b.Visitors {
		if err := v.VisitSystem(s); err != nil {
			return err
		}
	}
	return nil
}

func markFieldList(fl core.FieldList) {
	for i, f := range fl {
		f.Last = i == len(fl)-1
	}
}

func markType(t *core.Type) {
	markFieldList(t.Fields)
}

func markFunc(f *core.Func) {
	markFieldList(f.Params)
	markFieldList(f.Return)
}

func markModule(m *core.Module) {
	for _, t := range m.Types {
		markType(t)
	}
	for _, f := range m.Funcs {
		markFunc(f)
	}
}

func (b *Builder) VisitModule(module *core.Module) error {
	if _, seen := b.seen[module.Name]; seen {
		return nil
	}
	b.seen[module.Name] = true
	markModule(module)
	for _, v := range b.Visitors {
		if err := v.VisitModule(module); err != nil {
			return err
		}
	}
	return nil
}

func NewBuilder(visitors ...core.Visitor) core.Visitor {
	return &Builder{Visitors: visitors, seen: make(map[string]bool)}
}
