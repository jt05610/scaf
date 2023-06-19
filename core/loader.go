package core

import (
	"bytes"
	"embed"
	"fmt"
	"github.com/gertd/go-pluralize"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"io/fs"
	"path/filepath"
	"strings"
	"text/template"
)

// TypeCrawler visits all types and generates functions to translate inputs into the correct type
type TypeCrawler struct {
	seen  map[*Type]string
	l     *Language
	Types []*Type
	tpl   *template.Template
}

// VisitField visits a field and generates functions to translate inputs into the correct type
func (c *TypeCrawler) VisitField(f *Field) error {
	if !f.Type.IsPrimitive() {
		if err := c.VisitType(f.Type); err != nil {
			return err
		}
	}
	return nil
}

// VisitType visits a type and generates functions to translate inputs into the correct type
func (c *TypeCrawler) VisitType(t *Type) error {
	if _, seen := c.seen[t]; seen {
		return nil
	}
	var buf bytes.Buffer
	err := c.tpl.Execute(&buf, t)
	if err != nil {
		return err
	}

	c.seen[t] = buf.String()
	for _, f := range t.Fields {
		if err := c.VisitField(f); err != nil {
			return err
		}
	}

	return nil
}

func (c *TypeCrawler) Crawl() error {
	for _, t := range c.Types {
		if err := c.VisitType(t); err != nil {
			return err
		}
	}
	return nil
}

func (c *TypeCrawler) String() string {
	var buf bytes.Buffer
	for _, t := range c.Types {
		buf.WriteString(c.seen[t])
	}
	return buf.String()
}

func NewCrawler(l *Language, m *API) *TypeCrawler {
	tpl := template.Must(template.New("translateFuncs").Funcs(funcMap(l)).Parse(fmt.Sprintf(`
func New{{.Name}}(input *%s.{{.Name}}Input) *%s.{{.Name}} {
	ret := &%s.{{.Name}}{
	{{- range .Fields}}
		{{- if .Type.IsArray}}
	{{.Name}}: make([]*%s.{{.Type.Name}}, len(input.{{.Name}})),
		{{- else}}
	{{.Name}}: {{if .Type.IsPrimitive}}input.{{.Name}}{{else}}New{{typeTrans .}}(input.{{.Name}}){{end}},
		{{- end}}
	{{- end}}
	}

	{{- range .Fields}}
		{{- if .Type.IsArray}}
	for _, i := range input.{{.Name}} {
			{{- if .Type.IsPrimitive}}
		ret.{{.Name}} = append(ret.{{.Name}}, i)
			{{- else}}
		ret.{{.Name}} = append(ret.{{.Name}}, New{{.Type.Name}}(i))
			{{- end}}
	}
		{{- end}}
	{{- end}}
	return ret
}
`, m.Name, m.Name, m.Name, m.Name)))
	return &TypeCrawler{
		seen:  make(map[*Type]string),
		Types: m.Types,
		l:     l,
		tpl:   tpl,
	}
}

func funcMap(l *Language) template.FuncMap {
	base := template.FuncMap{
		"add": func(i, j int) int {
			return i + j
		},
		"typeTrans": func(f Field) string {
			return f.TypeString(l)
		},
		"inputTrans": func(f Field) string {
			return f.InputString(l)
		},
		"pluralize": pluralize.NewClient().Plural,
		"lower":     cases.Lower(language.Und).String,
		"upper":     cases.Upper(language.Und).String,
		"pascal":    cases.Title(language.Und).String,
		"translateFuncs": func(m *API) string {
			c := NewCrawler(l, m)
			if err := c.Crawl(); err != nil {
				panic(err)
			}
			return c.String()
		},
	}
	base["service"] = func(api *API) string {
		t, err := template.New("service").Funcs(base).Parse(api.Language.Service)
		if err != nil {
			panic(err)
		}
		var b strings.Builder
		err = t.Execute(&b, api)
		return b.String()
	}
	return base
}

type Entry struct {
	Path     *template.Template
	Template *template.Template
	Children []*Entry
}

type TemplateLoader struct {
	lang    *Language
	entries map[string]*Entry
	seen    map[string]bool
	fs      *embed.FS
	tplDir  string
	parent  string
}

func (l *TemplateLoader) pathName(path string) *template.Template {
	fn := strings.TrimSuffix(path, ".gotpl")
	modPath := filepath.Join(l.tplDir, "module")
	sysPath := filepath.Join(l.tplDir, "system")
	if strings.Contains(path, modPath) {
		fn = strings.ReplaceAll(fn, modPath, "{{.Name}}")
		fn = filepath.Join(l.parent, fn)
	} else if strings.Contains(path, sysPath) {
		fn = strings.ReplaceAll(fn, sysPath, l.parent)
	}
	return template.Must(template.New(path).Funcs(funcMap(l.lang)).Parse(fn))
}

func (l *TemplateLoader) loadDir(e *Entry, path string) error {
	e.Path = l.pathName(path)
	e.Children = make([]*Entry, 0)
	entries, err := l.fs.ReadDir(path)
	if err != nil {
		return err
	}
	for _, de := range entries {
		ep, err := l.loadEntry(path, de)
		if err != nil {
			return err
		}
		if ep != nil {
			e.Children = append(e.Children, ep)
		}
	}
	return nil
}

func (l *TemplateLoader) loadFile(e *Entry, path string) (err error) {
	f, err := l.fs.ReadFile(path)
	if err != nil {
		return err
	}
	e.Path = l.pathName(path)
	e.Template, err = template.New(path).Funcs(funcMap(l.lang)).Parse(string(f))
	if err != nil {
		return err
	}
	return nil
}

func (l *TemplateLoader) loadEntry(parent string, ff fs.DirEntry) (e *Entry, err error) {
	e = &Entry{}
	path := filepath.Join(parent, ff.Name())
	if _, seen := l.seen[path]; seen {
		return nil, nil
	}
	l.seen[path] = true
	if ff.IsDir() {
		err = l.loadDir(e, path)
	} else {
		err = l.loadFile(e, path)
	}
	if err != nil {
		return nil, err
	}
	return e, nil
}

func (l *TemplateLoader) Load() error {
	par, err := l.fs.ReadDir(".")
	if err != nil {
		return err
	}
	l.tplDir = filepath.Join(par[0].Name(), l.lang.Name)
	content, err := l.fs.ReadDir(l.tplDir)
	if err != nil {
		return err
	}
	for _, ff := range content {
		e, err := l.loadEntry(l.tplDir, ff)
		if err != nil {
			return err
		}
		l.entries[strings.TrimPrefix(ff.Name(), filepath.Join("template", l.lang.Name))] = e
	}
	return nil
}

func (l *TemplateLoader) Module() []*Entry {
	return l.entries["module"].Children
}

func (l *TemplateLoader) System() []*Entry {
	e, ok := l.entries["system"]
	if !ok {
		return nil
	}
	return e.Children
}

func NewLoader(parent string, lang *Language) *TemplateLoader {
	return &TemplateLoader{
		entries: make(map[string]*Entry, 0),
		seen:    make(map[string]bool),
		fs:      lang.FS(),
		parent:  parent,
		lang:    lang,
	}
}
