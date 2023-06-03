package system

import (
	"github.com/jt05610/scaf/caddy"
	"os/user"
	"strings"
	"time"
)

type Kind string

const (
	Device   Kind = "device"
	Software Kind = "software"
	Library  Kind = "library"
)

// System is the top-level struct for a system.
type System struct {
	Author  string `yaml:"author"`
	Date    string `yaml:"date"`
	Name    string `prompt:"What is the name of this system?" default:"System"`
	Kind    Kind   `prompt:"What kind of system is this?" options:"device,software,library" default:"device"`
	Modules []*Module
}

// FillMetadata fills in the metadata for the system.
func (s *System) FillMetadata() error {
	u, err := user.Current()
	if err != nil {
		return err
	}
	s.Author = u.Name
	s.Date = time.Now().Format("02 Jan 2006")
	return nil
}

func (s *System) AddModule(m *Module) error {
	if s.Modules == nil {
		s.Modules = make([]*Module, 0)
	}
	s.Modules = append(s.Modules, m)
	cf := s.Caddyfile()
	m.Addr = cf.Servers[len(s.Modules)-1].Addr
	m.Port = cf.Servers[len(s.Modules)-1].Port
	return nil
}

// Caddyfile returns a Caddyfile for the system.
func (s *System) Caddyfile() *caddy.Caddyfile {
	cf := caddy.NewCaddyfile(&caddy.Options{
		UIPortStart:  3000,
		APIPortStart: 8000,
		PortTimeout:  time.Duration(10) * time.Millisecond,
	}, strings.ToLower(s.Name)+".bot")
	for _, m := range s.Modules {
		var kind caddy.ServerKind
		if m.HasUi {
			kind = caddy.UI
		} else {
			kind = caddy.API
		}
		cf.AddServer(&caddy.Server{
			Kind: kind,
			Addr: "localhost",
			Path: "/" + m.Name,
		})
	}
	return cf
}
