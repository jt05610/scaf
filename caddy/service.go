package caddy

import (
	"fmt"
	"github.com/jt0610/scaf/codegen"
	"github.com/jt0610/scaf/service"
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
	Addr      string
	Servers   []*Server
	APIs      []*Server
	Frontends []*Server
	curPort   map[ServerKind]int
	portTO    time.Duration
}

func (c *Caddyfile) portIsOpen(port int) bool {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("localhost:%d", port), c.portTO)
	if err != nil {
		if _, ok := err.(*net.OpError); ok {
			return true
		}
	}
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
		if c.portIsOpen(i) {
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
		c.assignPort(s)
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
		portTO:    opt.PortTimeout,
		APIs:      make([]*Server, 0),
		Frontends: make([]*Server, 0),
	}
}

type renderer struct {
	parDir string
}

func (r *renderer) Load(rdr io.Reader) (*Caddyfile, error) {
	t, err := codegen.Load("caddyfile.gotpl")
	if err != nil {
		return nil, err
	}
	f := &Caddyfile{}
}

func (r *renderer) Suffix() string {
	return "Caddyfile"
}

func (r *renderer) Flush(w io.Writer, f *Caddyfile) error {
	for _, s := range f.Servers {
		if s.Kind == API {
			f.APIs = append(f.APIs, s)
		} else {
			f.Frontends = append(f.Frontends, s)
		}
	}

	t, err := codegen.Load("caddyfile.gotpl")
	if err != nil {
		return err
	}
	return t.Execute(w, f)
}

func Service(parDir string) service.Service[*Caddyfile] {
	return &renderer{parDir: parDir}
}
