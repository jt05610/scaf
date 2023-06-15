package core

import (
	"embed"
	"github.com/gertd/go-pluralize"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"io/fs"
	"path/filepath"
	"strings"
	"text/template"
)

func funcMap(l Language) template.FuncMap {
	return template.FuncMap{
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
	}
}

type Entry struct {
	Path     *template.Template
	Template *template.Template
	Children []*Entry
}

type TemplateLoader struct {
	lang    Language
	entries map[string]*Entry
	seen    map[string]bool
	fs      *embed.FS
	tplDir  string
	parent  string
}

func (l *TemplateLoader) pathName(path string) *template.Template {
	fn := strings.TrimSuffix(path, ".gotpl")
	if strings.Contains(path, "template/module") {
		fn = strings.ReplaceAll(fn, "template/module", "{{.Name}}")
		fn = filepath.Join(l.parent, fn)
	} else if strings.Contains(path, "template/system") {
		fn = strings.ReplaceAll(fn, "template/system", l.parent)
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
	l.tplDir = par[0].Name()
	content, err := l.fs.ReadDir(l.tplDir)
	if err != nil {
		return err
	}
	for _, ff := range content {
		e, err := l.loadEntry(l.tplDir, ff)
		if err != nil {
			return err
		}
		l.entries[strings.TrimPrefix(ff.Name(), "template/")] = e
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

func NewLoader(parent string, lang Language) *TemplateLoader {
	return &TemplateLoader{
		entries: make(map[string]*Entry, 0),
		seen:    make(map[string]bool),
		fs:      lang.FS(),
		parent:  parent,
		lang:    lang,
	}
}
