package caddy

import (
	"github.com/jt0610/scaf/file"
	"io"
)

type Server struct {
	Path string
	Port int
}

type Caddyfile struct {
	APIs []Server
	UIs  []Server
}

type Service struct{}

func (s Service) Load(r io.Reader) (*Caddyfile, error) {
	//TODO implement me
	panic("implement me")
}

func (s Service) Flush(w io.Writer, data *Caddyfile) error {
	//TODO implement me
	panic("implement me")
}

func New() file.Service[*Caddyfile] {
	return Service{}
}
