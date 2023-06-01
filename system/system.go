package system

type Kind string

const (
	Device   Kind = "device"
	Software Kind = "software"
	Library  Kind = "library"
)

type Language string

type System struct {
	Name string
	Kind Kind
}

func NewSystem(kind Kind, name string) *System {
	return &System{
		Name: name,
		Kind: kind,
	}
}
