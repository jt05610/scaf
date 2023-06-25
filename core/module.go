package core

type BaseType string

type Type interface {
	IsPrimitive() bool
	String() string
}

func (t BaseType) IsPrimitive() bool {
	return true
}

func (t BaseType) String() string {
	return string(t)
}

const (
	Int    BaseType = "int"
	Float  BaseType = "float"
	String BaseType = "string"
	ID     BaseType = "ID"
	Bool   BaseType = "bool"
)

func External(name, kind, pluralKind string) *Field {
	return &Field{
		Name: name,
		Type: &Model{
			IsExternal: true,
			Name:       kind,
			Plural:     pluralKind,
			Fields: []*Field{
				{
					Name:       "ID",
					Type:       ID,
					IsExternal: true,
				},
			},
		},
		IsExternal: true,
	}
}

type Model struct {
	IsExternal  bool
	Name        string
	Plural      string
	Description string
	Query       bool
	Create      bool
	Update      bool
	Delete      bool
	Fields      []*Field
}

func (m *Model) IsPrimitive() bool {
	return false
}

func (m *Model) String() string {
	return m.Name
}

// Field represents a field on a model
type Field struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Type        Type   `yaml:"type"`
	Required    bool   `yaml:"required,omitempty"`
	Last        bool   `yaml:"-"`
	IsArray     bool   `yaml:"IsArray,omitempty"`
	Query       bool   `yaml:"query,omitempty"`
	Create      bool   `yaml:"create,omitempty"`
	Update      bool   `yaml:"update,omitempty"`
	Delete      bool   `yaml:"delete,omitempty"`
	Subscribe   bool   `yaml:"subscribe,omitempty"`
	IsExternal  bool   `yaml:"is_external,omitempty"`
}

type Func struct {
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	UserCode    string   `yaml:"user_code"`
	Return      []*Field `yaml:"return"`
	Params      []*Field `yaml:"params"`
}

type PortMap struct {
	UI  int
	GQL int
	RPC int
}

func (p *PortMap) Add(n int) *PortMap {
	return &PortMap{
		UI:  p.UI + n,
		GQL: p.GQL + n,
		RPC: p.RPC + n,
	}
}

type Module struct {
	*MetaData
	Version int
	apis    []*API
}

func (m *Module) Meta() *MetaData {
	return m.MetaData
}

func NewModule(name, author, date string) *Module {
	return &Module{
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
		apis: make([]*API, 0),
	}
}

func (m *Module) APIs() []*API {
	return m.apis
}

func (m *Module) AddAPI(a *API) {
	m.apis = append(m.apis, a)
	m.Version = len(m.apis)
	a.Version = m.Version
}
