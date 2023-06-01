package system

// Parameter represents a parameter of a function.
type Parameter struct {
	// Name is the identifier of the parameter.
	Name string `prompt:"What is the name of this parameter?" default:"Default Parameter"`

	// Desc is a brief description of what this parameter represents.
	Desc string `prompt:"Provide a brief description for this parameter."`

	// Type specifies the data type of the parameter.
	Type string `prompt:"What is the data type of this parameter?" options:"string,int,bool"`

	// Query specifies whether this parameter can be queried.
	Query bool `prompt:"Should this parameter be stored and viewed later?"`

	// Mutate specifies whether this parameter can be changed.
	Mutate bool `prompt:"Can this parameter be changed?"`

	// Subscribe specifies whether updates to this parameter can be subscribed to.
	Subscribe bool `prompt:"Can updates to this parameter be subscribed to?"`
}
