package codegen

import (
	"bytes"
	"embed"
	"github.com/gertd/go-pluralize"
	"github.com/jt05610/scaf/core"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

//go:embed template
var templates embed.FS

var Commands = `
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
protoc --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. v1/{{.Name}}.proto
go mod tidy
go generate ./...
mkcert {{.Name}}.local localhost 127.0.0.1 ::1
go fmt ./...
`
var funcMap = template.FuncMap{
	"add": func(i, j int) int {
		return i + j
	},
	"typeTrans": func(l core.Language, f core.Field) string {
		return f.TypeString(l)
	},
	"inputTrans": func(l core.Language, f core.Field) string {
		return f.InputString(l)
	},
	"pluralize": pluralize.NewClient().Plural,
	"lower":     cases.Lower(language.Und).String,
	"upper":     cases.Upper(language.Und).String,
	"pascal":    cases.Title(language.Und).String,
}

type Entry struct {
	Path     *template.Template
	Template *template.Template
	Children []*Entry
}

type generator struct {
	entries []*Entry
	tplDir  string
	seen    map[fs.DirEntry]bool
	parent  string
	fs      embed.FS
}

func (g *generator) nameTemplate(n string) *template.Template {
	fn := strings.ReplaceAll(strings.TrimSuffix(n, ".gotpl"), g.tplDir, "{{.Name}}")
	fn = filepath.Join(g.parent, fn)
	return template.Must(template.New(n).Funcs(funcMap).Parse(fn))
}

func (g *generator) load(parent string, ff fs.DirEntry) *Entry {
	if _, seen := g.seen[ff]; seen {
		return nil
	}
	g.seen[ff] = true
	path := filepath.Join(parent, ff.Name())
	entry := &Entry{Path: g.nameTemplate(path), Children: make([]*Entry, 0)}
	if ff.IsDir() {
		dirEntries, err := g.fs.ReadDir(path)
		if err != nil {
			panic(err)
		}
		for _, de := range dirEntries {
			ep := g.load(path, de)
			if ep != nil {
				entry.Children = append(entry.Children, ep)
			}
		}
	} else {
		f, err := g.fs.ReadFile(path)
		if err != nil {
			panic(err)
		}
		entry.Template, err = template.New(path).Funcs(funcMap).Parse(string(f))
		if err != nil {
			panic(err)
		}
	}
	return entry
}

func (g *generator) gen(m *core.Module, e *Entry) error {
	var pathBuffer bytes.Buffer
	err := e.Path.Execute(&pathBuffer, m)
	if err != nil {
		return err
	}
	path := pathBuffer.String()

	if e.Template != nil {
		wr, err := os.Create(path)
		if err != nil {
			return err
		}
		defer func() {
			_ = wr.Close()
		}()
		return e.Template.Execute(wr, m)
	}
	err = os.MkdirAll(path, 0755)
	if err != nil {
		return err
	}
	for _, c := range e.Children {
		if err := g.gen(m, c); err != nil {
			return err
		}
	}
	return nil
}

func (g *generator) Visit(m *core.Module) error {
	for _, e := range g.entries {
		if err := g.gen(m, e); err != nil {
			return err
		}
	}
	return nil
}

func Generator(parent string) core.Visitor {
	par, err := templates.ReadDir(".")
	if err != nil {
		panic(err)
	}

	g := &generator{
		entries: make([]*Entry, 0),
		seen:    make(map[fs.DirEntry]bool),
		fs:      templates,
		tplDir:  par[0].Name(),
		parent:  parent,
	}
	if err != nil {
		panic(err)
	}
	entries, err := g.fs.ReadDir(g.tplDir)
	if err != nil {
		panic(err)
	}
	for _, entry := range entries {
		g.entries = append(g.entries, g.load(g.tplDir, entry))
	}
	return g
}
