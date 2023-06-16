package codegen

import (
	"bytes"
	"github.com/jt05610/scaf/context"
	"github.com/jt05610/scaf/core"
	"go.uber.org/zap"
	"os"
	"strings"
)

type Generator struct {
	loader *core.TemplateLoader
}

func (g *Generator) gen(item interface{}, e *core.Entry) error {
	var pathBuffer bytes.Buffer
	err := e.Path.Execute(&pathBuffer, item)
	if err != nil {
		return err
	}
	path := pathBuffer.String()

	if e.Template != nil {
		split := strings.Split(path, "/")
		parent := strings.Join(split[:len(split)-1], "/")
		err = os.MkdirAll(parent, 0755)
		if err != nil {
			return err
		}
		wr, err := os.Create(path)
		if err != nil {
			return err
		}
		defer func() {
			_ = wr.Close()
		}()
		return e.Template.Execute(wr, item)
	}

	err = os.MkdirAll(path, 0755)
	if err != nil {
		return err
	}
	for _, c := range e.Children {
		if err := g.gen(item, c); err != nil {
			return err
		}
	}
	return nil
}

func (g *Generator) VisitModule(ctx context.Context, m *core.Module) error {
	ctx.Logger.Debug("Generating", zap.String("module", m.Name))
	for _, e := range g.loader.Module() {
		ctx.Logger.Debug("Generating", zap.String("template", e.Path.Name()))
		if err := g.gen(m, e); err != nil {
			ctx.Logger.Error("Error generating", zap.String("template", e.Path.Name()), zap.Error(err))
			return err
		}
	}
	return nil
}

func (g *Generator) VisitSystem(ctx context.Context, s *core.System) error {
	ctx.Logger.Debug("Generating", zap.String("system", s.Name))
	err := g.loader.Load()
	if err != nil {
		return err
	}
	for _, e := range g.loader.System() {
		if err := g.gen(s, e); err != nil {
			ctx.Logger.Error("Error generating", zap.String("template", e.Path.Name()), zap.Error(err))
			return err
		}
	}
	for _, m := range s.Modules {
		ctx.Logger.Debug("Generating", zap.String("module", m.Name))
		if err := g.VisitModule(ctx, m); err != nil {
			return err
		}
	}
	return nil
}

func New(parent string, lang core.Language) *Generator {
	loader := core.NewLoader(parent, lang)
	return &Generator{loader: loader}
}
