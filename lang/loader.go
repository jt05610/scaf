package lang

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"github.com/gertd/go-pluralize"
	"github.com/jt05610/scaf/core"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"io/fs"
	"path/filepath"
	"strings"
	"text/template"
)

// TypeCrawler visits all types and generates functions to translate inputs into the correct type
type TypeCrawler struct {
	seen  map[*core.Model]string
	l     *Language
	Types []*core.Model
	tpl   *template.Template
}

// VisitType visits a type and generates functions to translate inputs into the correct type
func (c *TypeCrawler) VisitType(t core.Type) error {
	if t.IsPrimitive() {
		return nil
	}
	m := t.(*core.Model)
	if _, seen := c.seen[m]; seen {
		return nil
	}

	var buf bytes.Buffer
	err := c.tpl.Execute(&buf, m)
	if err != nil {
		return err
	}

	c.seen[m] = buf.String()
	for _, f := range m.Fields {
		if err := c.VisitType(f.Type); err != nil {
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
	for _, str := range c.seen {
		buf.WriteString(str)
	}
	return buf.String()
}

func NewCrawler(l *Language, m *core.API) *TypeCrawler {
	tpl := template.New("translateFuncs").Funcs(funcMap(l))
	if l.Name == "go" {
		tpl = template.Must(tpl.Parse(fmt.Sprintf(`
func New{{.Name}}(input *%s.{{.Name}}Input) *%s.{{.Name}} {
	ret := &%s.{{.Name}}{
	{{- range .Fields}}
		{{- if .IsArray}}
	{{.Name}}: make([]*%s.{{.Type.Name}}, len(input.{{.Name}})),
		{{- else}}
	{{.Name}}: {{if .Type.IsPrimitive}}input.{{.Name}}{{else}}New{{typeTrans .}}(input.{{.Name}}){{end}},
		{{- end}}
	{{- end}}
	}

	{{- range .Fields}}
		{{- if .IsArray}}
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
	} else if l.Name == "gql" {
		tpl = template.Must(tpl.Parse(fmt.Sprintf(`
{{- if .IsExternal}}
extend type {{.Name}} @key(fields: "id") {
    id: ID! @external
}

{{- else}}
type {{.Name}} {{if .Query}}@key(fields: "id"){{end}}{
	{{if .Query}}id: ID!{{end}}
    {{- range .Fields}}
    {{lower .Name}}: {{typeTrans .}}
    {{- end}}
}

input {{.Name}}Input {
    {{- range .Fields}}
        {{- if not .IsExternal}}
            {{- if .Create}}
    {{lower .Name}}: {{inputTrans .}}{{if.Required}}!{{end }}
            {{- end}}
        {{- else}}
	{{lower .Name}}: ID
        {{- end}}
    {{- end}}
}
{{- end}}
`)))
	} else if l.Name == "sql" {
		tpl = template.Must(tpl.Parse(fmt.Sprintf(`
{{- $model := .}}
{{- if not .IsExternal}}
create table {{lower .Plural}} (
	id bigserial primary key,
	{{- range $i, $e := .Fields}}
		{{- if not .IsArray -}}
	{{if gt $i 0}}, {{end}}
	"{{fLower .Name}}" {{typeTrans .}}
		{{- else}}
		{{- end}}
	{{- end}}
);

	{{- range .Fields}}
		{{- if .IsArray}}
create table {{lower .Name}} (
	id bigserial primary key,
	"{{lower .Name}}_id" bigint 
	constraint fk_{{lower $model.Name}} foreign key({{lower .Name}}_id) references {{lower $model.Plural}}(id)
);
		{{- end}}
	{{- end}}
{{- end}}
`)))
	}
	return &TypeCrawler{
		seen:  make(map[*core.Model]string),
		Types: m.Models,
		l:     l,
		tpl:   tpl,
	}
}

func funcMap(l *Language) template.FuncMap {
	base := template.FuncMap{
		"add": func(i, j int) int {
			return i + j
		},
		"typeTrans": func(f core.Field) string {
			return l.TypeDecl(f)
		},
		"inputTrans": func(f core.Field) string {
			return l.InputDecl(f)
		},
		"pluralize": pluralize.NewClient().Plural,
		"lower":     cases.Lower(language.Und).String,
		"fLower": func(s string) string {
			return strings.ToLower(s[:1]) + s[1:]
		},
		"upper":  cases.Upper(language.Und).String,
		"pascal": cases.Title(language.Und).String,
		"translateFuncs": func(m *core.API) string {
			c := NewCrawler(l, m)
			if err := c.Crawl(); err != nil {
				panic(err)
			}
			return c.String()
		},
		"translateMod": func(m *core.Module) string {
			c := NewCrawler(l, m.APIs()[m.Version-1])
			if err := c.Crawl(); err != nil {
				panic(err)
			}
			return c.String()
		},
		"makeScripts": func(m *core.API) string {
			var mm map[string]string
			switch m.Language {
			case core.Go:
				mm = Go("").Scripts.Map
			case core.Python:
				mm = Python("").Scripts.Map
			default:
				return ""
			}
			var buf bytes.Buffer
			err := json.NewEncoder(&buf).Encode(mm)
			if err != nil {
				panic(err)
			}
			t, err := template.New("service").Parse(buf.String())
			if err != nil {
				panic(err)
			}
			var out bytes.Buffer
			err = t.Execute(&out, m)
			if err != nil {
				panic(err)
			}
			return out.String()
		},
	}
	base["service"] = func(api *core.API) string {
		var t *template.Template
		var err error
		if api.Language == core.Go {
			t, err = template.New("service").Funcs(base).Parse(Go("").Service)
		} else if api.Language == core.Python {
			t, err = template.New("service").Funcs(base).Parse(Python("").Service)
		}
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
