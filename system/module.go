package system

type UIType string

// Module represents a specific component or functionality of the System.
type Module struct {
	// Name is the identifier of the module.
	Name string `prompt:"What is the name of this module?" default:"Module"`

	// Name is the identifier of the module.
	Desc string `prompt:"What does this module do? What problem does it solve?" default:"Module"`

	// HasUI indicates whether this module has a UI.
	HasUI bool `prompt:"Does this module have a UI?" default:"false"`

	// Comm is the way this module communicates with other modules.
	Communication *Comm

	// Resources are the associated data or files needed for this module.
	Resources []*Resource

	// Externals are any external dependencies or modules needed by this module.
	Externals []*External
}
