package system

type Kind string

const (
	Device   Kind = "device"
	Software Kind = "software"
	Library  Kind = "library"
)

type Language string

type System struct {
	Name    string `prompt:"What is the name of this system?" default:"System"`
	Kind    Kind   `prompt:"What kind of system is this?" options:"device,software,library" default:"device"`
	Modules []*Module
}

func (s *System) FillMetadata(author string) {

}
