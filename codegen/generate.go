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
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
)

//go:embed template
var templates embed.FS

//go:embed system
var sysTemplates embed.FS

var Commands = `
mkcert {{.Name}}.local localhost 127.0.0.1 ::1
mv {{.Name}}.local+3.pem {{.Name}}.local+3-key.pem ./cmd/.secrets
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
protoc --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. v1/{{.Name}}.proto
go mod tidy
go generate ./...
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

type Generator struct {
	entries []*Entry
	tplDir  string
	seen    map[fs.DirEntry]bool
	parent  string
	fs      embed.FS
}

func (g *Generator) nameTemplate(n string, fsName string, newValue string) *template.Template {
	fn := strings.ReplaceAll(strings.TrimSuffix(n, ".gotpl"), fsName, newValue)
	fn = filepath.Join(g.parent, fn)
	return template.Must(template.New(n).Funcs(funcMap).Parse(fn))
}

func (g *Generator) load(parent string, fsName string, fs embed.FS, ff fs.DirEntry) *Entry {
	if _, seen := g.seen[ff]; seen {
		return nil
	}
	g.seen[ff] = true
	path := filepath.Join(parent, ff.Name())
	var entry *Entry
	if fsName == "system" {
		entry = &Entry{Path: g.nameTemplate(path, fsName, "."), Children: make([]*Entry, 0)}
	} else {

		entry = &Entry{Path: g.nameTemplate(path, fsName, "{{.Name}}"), Children: make([]*Entry, 0)}
	}
	if ff.IsDir() {
		dirEntries, err := fs.ReadDir(path)
		if err != nil {
			panic(err)
		}
		for _, de := range dirEntries {
			ep := g.load(path, fsName, fs, de)
			if ep != nil {
				entry.Children = append(entry.Children, ep)
			}
		}
	} else {
		f, err := fs.ReadFile(path)
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

func (g *Generator) gen(m *core.Module, e *Entry) error {
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

func (g *Generator) sysGen(s *core.System, e *Entry) error {
	var pathBuffer bytes.Buffer
	err := e.Path.Execute(&pathBuffer, s)
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
		return e.Template.Execute(wr, s)
	}
	err = os.MkdirAll(path, 0755)
	if err != nil {
		return err
	}
	for _, c := range e.Children {
		if err := g.sysGen(s, c); err != nil {
			return err
		}
	}
	return nil
}

func (g *Generator) Visit(m *core.Module) error {
	for _, e := range g.entries {
		if err := g.gen(m, e); err != nil {
			return err
		}
	}
	return nil
}

func (g *Generator) Init(s *core.System) error {
	par, err := sysTemplates.ReadDir(".")
	if err != nil {
		return err
	}
	entries, err := sysTemplates.ReadDir(par[0].Name())
	if err != nil {
		return err
	}
	sysEntries := make([]*Entry, 0)
	for _, entry := range entries {
		sysEntries = append(sysEntries, g.load(par[0].Name(), par[0].Name(), sysTemplates, entry))
	}
	for _, e := range sysEntries {
		if err := g.sysGen(s, e); err != nil {
			return err
		}
	}
	cmd := exec.Command("mv", "gitignore", ".gitignore")
	cmd.Dir = filepath.Join(g.parent)
	err = cmd.Run()
	if err != nil {
		return err
	}

	cmd = exec.Command("git", "init", "&&", "git", "add", "*")
	cmd.Dir = filepath.Join(g.parent)
	err = cmd.Run()
	if err != nil {
		return err
	}

	cmd = exec.Command("npm", "install")
	cmd.Dir = filepath.Join(g.parent)
	err = cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func New(parent string) *Generator {
	par, err := templates.ReadDir(".")
	if err != nil {
		panic(err)
	}

	g := &Generator{
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
		g.entries = append(g.entries, g.load(g.tplDir, g.tplDir, templates, entry))
	}
	return g
}