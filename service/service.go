package service

import "io"

type Renderer[T any] interface {
	Flush(w io.Writer, t T) error
}

type Loader[T any] interface {
	Load(r io.Reader) (T, error)
}

type Service[T any] interface {
	Renderer[T]
	Loader[T]
	Suffix() string
}
