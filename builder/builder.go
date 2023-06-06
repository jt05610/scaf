package builder

import "github.com/jt05610/scaf/core"

type Builder struct {
	Visitors []core.Visitor
	seen     map[string]bool
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

func (b *Builder) Visit(module *core.Module) core.Visitor {
	if _, seen := b.seen[module.Name]; seen {
		return nil
	}
	b.seen[module.Name] = true
	markModule(module)
	for _, v := range b.Visitors {
		v.Visit(module)
	}
	return b
}

func NewBuilder(visitors ...core.Visitor) *Builder {
	v := make([]core.Visitor, 1)
	v[0] = core.NewChecker()
	v = append(v, visitors...)
	return &Builder{Visitors: v, seen: make(map[string]bool)}
}
