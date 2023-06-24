package builder

import (
	"github.com/jt05610/scaf/context"
	"github.com/jt05610/scaf/core"
	"go.uber.org/zap"
)

type Builder struct {
	Visitors []core.Visitor
	seen     map[string]bool
}

func (b *Builder) VisitSystem(ctx context.Context, s *core.System) error {
	ctx.Logger.Debug("Building", zap.String("system", s.Name))
	for _, v := range b.Visitors {
		if err := v.VisitSystem(ctx, s); err != nil {
			return err
		}
	}
	return nil
}

func markFieldList(fl []*core.Field) {
	fl[len(fl)-1].Last = true
}

func markType(t *core.Model) {
	markFieldList(t.Fields)
}

func markFunc(f *core.Func) {
	markFieldList(f.Params)
	markFieldList(f.Return)
}

func markModule(m *core.Module) {
	for _, api := range m.APIs() {

		for _, t := range api.Models {
			markType(t)
		}
		for _, f := range api.Funcs {
			markFunc(f)
		}
	}
}

func (b *Builder) VisitModule(ctx context.Context, module *core.Module) error {
	ctx.Logger.Debug("Building", zap.String("module", module.Name))
	if _, seen := b.seen[module.Name]; seen {
		return nil
	}
	b.seen[module.Name] = true
	markModule(module)
	for _, v := range b.Visitors {
		if err := v.VisitModule(ctx, module); err != nil {
			return err
		}
	}
	return nil
}

func NewBuilder(visitors ...core.Visitor) core.Visitor {
	return &Builder{Visitors: visitors, seen: make(map[string]bool)}
}
