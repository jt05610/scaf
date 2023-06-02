package system

import (
	"os/user"
	"time"
)

type Kind string

const (
	Device   Kind = "device"
	Software Kind = "software"
	Library  Kind = "library"
)

type Language string

type System struct {
	Author  string `yaml:"author"`
	Date    string `yaml:"date"`
	Name    string `prompt:"What is the name of this system?" default:"System"`
	Kind    Kind   `prompt:"What kind of system is this?" options:"device,software,library" default:"device"`
	Modules []*Module
}

func (s *System) FillMetadata() error {
	u, err := user.Current()
	if err != nil {
		return err
	}
	s.Author = u.Name
	s.Date = time.Now().Format("02 Jan 2006")
	return nil
}

func (s *System) AddModule(m *Module) {
	if s.Modules == nil {
		s.Modules = make([]*Module, 0)
	}
	s.Modules = append(s.Modules, m)
}

func (s *System) UpdateCaddy() {
}
