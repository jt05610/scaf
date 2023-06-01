package system

// Function represents an operation that a module can perform.
type Function struct {
	// Name is the identifier of the function.
	Name string `prompt:"What is the name of this function?" default:"My Function"`

	// Desc is a brief description of what this parameter represents.
	Desc string `prompt:"Provide a brief description for this function."`

	// Parameters holds the list of parameters associated with the function.
	Parameters []*Parameter
}
