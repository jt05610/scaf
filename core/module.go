package core

type BaseType string

const (
	IntType    BaseType = "int"
	FloatType  BaseType = "float"
	StringType BaseType = "string"
	BoolType   BaseType = "bool"
)

type Type struct {
	Name        string
	Plural      string
	Fields      []*Field
	IsArray     bool
	Query       bool
	Mutate      bool
	Subscribe   bool
	isPrimitive bool
}

func (t *Type) IsPrimitive() bool {
	return t.isPrimitive
}

var Int = &Type{Name: "int", isPrimitive: true}
var Float = &Type{Name: "float", isPrimitive: true}
var String = &Type{Name: "string", isPrimitive: true}
var Bool = &Type{Name: "bool", isPrimitive: true}

func Array(t *Type) *Type {
	return &Type{
		Name:      t.Name,
		Plural:    t.Plural,
		Fields:    t.Fields,
		IsArray:   true,
		Query:     t.Query,
		Mutate:    t.Mutate,
		Subscribe: t.Subscribe,
	}
}

func (t *Type) string(l *Language, isInput bool) (ret string) {
	var found bool
	if ret, found = l.MapType(BaseType(t.Name)); !found {
		ret = t.Name
		if isInput {
			ret += "Input"
		}
	}
	if t.IsArray {
		return l.MakeArray(ret)
	}
	return ret
}

func (t *Type) String(l *Language) string {
	return t.string(l, false)
}

func (t *Type) InputString(l *Language) string {
	return t.string(l, true)
}

var old = `

func makeArray(l Language, s string) string {
	switch l {
	case Clojure:
		return fmt.Sprintf("[]%s", s)
	case C:
		return fmt.Sprintf("%s *", s)
	case Cpp:
		return fmt.Sprintf("std::vector<%s>", s)
	case Rust:
		return fmt.Sprintf("Vec<%s>", s)
	default:
		return s
	}
}

var TypeMapping = map[Language]map[BaseType]string{
	C: {
		IntType:    "int",
		FloatType:  "float",
		StringType: "char*",
		BoolType:   "bool",
	},
	Cpp: {
		IntType:    "int",
		FloatType:  "float",
		StringType: "std::string",
		BoolType:   "bool",
	},
	Clojure: {
		IntType:    "integer",
		FloatType:  "float",
		StringType: "string",
		BoolType:   "boolean",
	},
	Rust: {
		IntType:    "i32",
		FloatType:  "f32",
		StringType: "String",
		BoolType:   "bool",
	},
}
`

type Field struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Required    bool   `yaml:"required,omitempty"`
	Type        *Type  `yaml:"type"`
	Last        bool   `yaml:"-"`
}

func (f *Field) IsPrimitive() bool {
	return f.Type.IsPrimitive()

}

func (f *Field) TypeString(l *Language) string {
	return f.Type.String(l)
}

func (f *Field) InputString(l *Language) string {
	return f.Type.InputString(l)
}

type Func struct {
	Name        string
	Description string
	Return      []*Field
	Params      []*Field
}

type API struct {
	Name     string    `yaml:"name"`
	PortMap  *PortMap  `yaml:"-"`
	HasSubs  bool      `yaml:"-"`
	Author   string    `yaml:"author"`
	Version  int       `yaml:"version"`
	Language *Language `yaml:"language"`
	Date     string    `yaml:"date"`
	Deps     []*Module `yaml:"deps"`
	Types    []*Type   `yaml:"types"`
	Funcs    []*Func   `yaml:"funcs"`
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
	API     map[string]*API
}

func (m *Module) Meta() *MetaData {
	return m.MetaData
}

func NewModule(name, author, date string, lang *Language) *Module {
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
		API: map[string]*API{
			"v1": {
				Version:  1,
				Language: lang,
				Date:     date,
				Deps:     make([]*Module, 0),
				Types:    make([]*Type, 0),
				Funcs:    make([]*Func, 0),
			},
		},
	}
}

func (m *Module) AddType(t *Type) {
	m.API["v1"].Types = append(m.API["v1"].Types, t)
}

func (m *Module) AddFunc(f *Func) {
	m.API["v1"].Funcs = append(m.API["v1"].Funcs, f)
}

func (m *Module) AddDep(d *Module) {
	m.API["v1"].Deps = append(m.API["v1"].Deps, d)
}
