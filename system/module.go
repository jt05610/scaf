package system

type UIType string

type Language string

const (
	Go         Language = "go"
	Python     Language = "python"
	Fortran    Language = "fortran"
	C          Language = "c"
	TypeScript Language = "typescript"
	Excel      Language = "excel"
)

// Field represents a field in a struct.
type Field struct {
	Name string
	Type string
}

// Module represents a specific component or functionality of the System.
type Module struct {
	// Author is the author of the module.
	Author string
	// Name is the identifier of the module.
	Name string `prompt:"What is the name of this module?" default:"Module"`

	// Desc is the description of the module.
	Desc string `prompt:"What does this module do? What problem does it solve?" default:"Module"`

	// Language is the programming language of this module.
	Language Language `prompt:"What language is this module written in?" options:"go,python,fortran,c" default:"go"`

	// Comm is the way this module communicates with other modules.
	Communication *Comm

	Functions []*Function

	HasUi bool `prompt:"Does this module have a UI?" default:"false"`

	// Externals are any external dependencies or modules needed by this module.
	Externals []string

	Fields []*Field

	Addr string
	Port int
}

type ModuleGenerator interface {
	// Generate generates the module.
	Generate(m *Module) error
}

func (m *Module) AddFunction(f *Function) {
	if m.Functions == nil {
		m.Functions = make([]*Function, 0)
	}
	m.Functions = append(m.Functions, f)
}
