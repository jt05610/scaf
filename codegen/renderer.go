package codegen

import "io"

type Renderer[T any] interface {
	Render(w io.Writer, t T) error
}
