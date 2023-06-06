package builder

import (
	"github.com/gertd/go-pluralize"
	"github.com/jt05610/scaf/core"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"text/template"
)

var funcMap = template.FuncMap{
	"add": func(i, j int) int {
		return i + 1
	},
	"typeTrans": func(l core.Language, f core.Field) string {
		return f.TypeString(l)
	},
	"pluralize": pluralize.NewClient().Plural,
	"lower":     cases.Lower(language.Und).String,
	"upper":     cases.Upper(language.Und).String,
	"pascal":    cases.Title(language.Und).String,
}

func NewTemplate(t string) *template.Template {
	return template.Must(
		template.New("proto").Funcs(funcMap).Parse(t),
	)
}

func FileTemplate(name func(m *core.Module) string, t string) *File {
	return &File{
		Name:     name,
		Template: NewTemplate(t),
	}
}
