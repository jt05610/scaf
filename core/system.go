package core

import "github.com/google/uuid"

type MetaData struct {
	ID          string   `yaml:"id" json:"_id,omitempty"`
	Rev         string   `yaml:"rev" json:"_rev,omitempty"`
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	Author      string   `yaml:"author"`
	Date        string   `yaml:"date"`
	PortMap     *PortMap `yaml:"ports"`
}

type Storable interface {
	Meta() *MetaData
}

type System struct {
	*MetaData
	Modules []*Module `yaml:"modules"`
}

func (s *System) Meta() *MetaData {
	return s.MetaData
}

func (s *System) FilterValue() string {
	return s.MetaData.Name
}

func (s *System) Title() string {
	return s.MetaData.Name
}

func (s *System) Description() string {
	return s.MetaData.Description
}

const (
	UIPort  = 4000
	GQLPort = 8080
	RPCPort = 5000
)

func id() string {
	u := uuid.New()
	return u.String()
}

func NewSystem(name, author, date string) *System {
	return &System{
		MetaData: &MetaData{
			ID:     id(),
			Name:   name,
			Author: author,
			Date:   date,
			PortMap: &PortMap{
				UI:  UIPort,
				GQL: GQLPort,
				RPC: RPCPort,
			},
		},
		Modules: make([]*Module, 0),
	}
}

func (s *System) AddModule(m *Module) error {
	m.PortMap = s.PortMap.Add(1)
	s.Modules = append(s.Modules, m)
	return nil
}
