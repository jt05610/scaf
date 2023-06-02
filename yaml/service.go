package yaml

import (
	"github.com/jt05610/scaf/service"
	"gopkg.in/yaml.v3"
	"io"
)

type Service[T any] struct{}

func (s *Service[T]) Suffix() string {
	return ".yaml"
}

func (s *Service[T]) Load(r io.Reader) (T, error) {
	b, err := io.ReadAll(r)
	if err != nil {
		var zero T
		return zero, err
	}
	var ret T
	return ret, yaml.Unmarshal(b, &ret)
}

func (s *Service[T]) Flush(w io.Writer, t T) error {
	bytes, err := yaml.Marshal(t)
	if err != nil {
		return err
	}
	_, err = w.Write(bytes)
	return err
}

func NewYAMLService[T any]() service.Service[T] {
	return &Service[T]{}
}
