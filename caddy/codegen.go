package caddy

import (
	"fmt"
	"github.com/jt0610/scaf/codegen"
	"io"
	"net"
	"time"
)

type ServerKind string

const (
	API ServerKind = "api"
	UI  ServerKind = "ui"
)

type Server struct {
	Kind ServerKind
	Path string
	Addr string
	Port int
}

type Caddyfile struct {
	Addr    string
	Servers []*Server
	curPort map[ServerKind]int
	portTO  time.Duration
}

func (c *Caddyfile) portIsOpen(port int) bool {
	conn, _ := net.DialTimeout("tcp", fmt.Sprintf("localhost:%d", port), c.portTO)
	if conn != nil {
		defer func(conn net.Conn) {
			err := conn.Close()
			if err != nil {
				panic(err)
			}
		}(conn)
		return false
	}
	return true
}

func (c *Caddyfile) nextOpenPort(start int, try int) (int, error) {
	for i := start; i < start+try; i++ {
		if !c.portIsOpen(i) {
			return i, nil
		}
	}
	return 0, fmt.Errorf("no open ports found in range")
}

func (c *Caddyfile) assignPort(s *Server) {
	s.Port = c.curPort[s.Kind]
	next, err := c.nextOpenPort(s.Port, 10)
	if err != nil {
		panic(err)
	}
	c.curPort[s.Kind] = next
}

func (c *Caddyfile) AddServer(s *Server) {
	if s.Addr == "" {
		s.Addr = "localhost"
	}
	if s.Port == 0 {
		a
	}
	c.Servers = append(c.Servers, s)
}

func NewCaddyfile(opt *codegen.Options, addr string) *Caddyfile {
	return &Caddyfile{
		Addr:    addr,
		Servers: make([]*Server, 0),
		curPort: map[ServerKind]int{
			API: opt.APIPortStart,
			UI:  opt.UIPortStart,
		},
	}
}

type renderer struct {
	parDir string
}

func (r renderer) Render(w io.Writer, t *Caddyfile) error {
}

func NewRenderer(parDir string) codegen.Renderer[*Caddyfile] {
	return &renderer{parDir: parDir}
}
