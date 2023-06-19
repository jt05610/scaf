package core

import "github.com/google/uuid"

type MetaData struct {
	ID          string   `yaml:"id" json:"_id,omitempty"`
	Rev         string   `yaml:"rev" json:"_rev,omitempty"`
	Name        string   `yaml:"name" json:"name,omitempty"`
	Description string   `yaml:"description" json:"description,omitempty"`
	Author      string   `yaml:"author" json:"author,omitempty"`
	Date        string   `yaml:"date" json:"date,omitempty"`
	PortMap     *PortMap `yaml:"ports" json:"port_map,omitempty"`
}

type Storable interface {
	Meta() *MetaData
}

type System struct {
	*MetaData
	Modules []*Module `yaml:"modules" json:"modules,omitempty"`
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

func NewSystem(name, description, author, date string) *System {
	return &System{
		MetaData: &MetaData{
			ID:          id(),
			Name:        name,
			Description: description,
			Author:      author,
			Date:        date,
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
