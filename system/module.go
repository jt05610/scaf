package system

type UIType string

type Language string

const (
	Go         Language = "go"
	Python     Language = "python"
	Fortran    Language = "fortran"
	C          Language = "c"
	TypeScript Language = "typescript"
)

// Module represents a specific component or functionality of the System.
type Module struct {
	// Name is the identifier of the module.
	Name string `prompt:"What is the name of this module?" default:"Module"`

	// Desc is the description of the module.
	Desc string `prompt:"What does this module do? What problem does it solve?" default:"Module"`

	// Language is the programming language of this module.
	Language Language `prompt:"What language is this module written in?" options:"go,python,fortran,c" default:"go"`

	// Comm is the way this module communicates with other modules.
	Communication *Comm

	Functions []*Function

	// Externals are any external dependencies or modules needed by this module.
	Externals []string
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
