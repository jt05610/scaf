package file

import "io"

type Service[T any] interface {
	Load(r io.Reader) (T, error)
	Flush(w io.Writer, data T) error
}
