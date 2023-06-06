package gql

import (
	"github.com/jt05610/core"
	"io"
)

type Printer struct {
	seen map[string]bool
	w    io.Writer
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

func (w *Printer) Visit(m *core.Module) error {
	if _, seen := w.seen[m.Name]; seen {
		return nil
	}
	w.seen[m.Name] = true
	markModule(m)
	t := gqlTemplate
	return t.Execute(w.w, m)
}

func NewPrinter(w io.Writer) *Printer {
	return &Printer{seen: make(map[string]bool), w: w}
}
