package caddy

import (
	"github.com/jt0610/scaf/codegen"
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

type renderer struct {
	parDir string
}

func (r renderer) Render(w io.Writer, t *Caddyfile) error {
	return nil
}

func NewRenderer(parDir string) codegen.Renderer[*Caddyfile] {
	return &renderer{parDir: parDir}
}
