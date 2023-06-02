package codegen

import (
	"html/template"
)

func Load(name string) (*template.Template, error) {
	return template.ParseFiles(name)
}
