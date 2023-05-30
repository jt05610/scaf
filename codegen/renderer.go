package codegen

import (
	"html/template"
	"io"
)

type Renderer[T any] interface {
	Render(w io.Writer, t T) error
}

func Load(name string) (*template.Template, error) {
	return template.ParseFiles(name)
}
