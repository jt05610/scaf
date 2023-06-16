package core

type MetaData struct {
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	Author      string   `yaml:"author"`
	Date        string   `yaml:"date"`
	PortMap     *PortMap `yaml:"ports"`
}

type System struct {
	*MetaData
	Modules []*Module `yaml:"modules"`
}

const (
	UIPort  = 4000
	GQLPort = 8080
	RPCPort = 5000
)

func NewSystem(name, author, date string) *System {
	return &System{
		MetaData: &MetaData{
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
