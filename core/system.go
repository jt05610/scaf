package core

type System struct {
	Name    string    `yaml:"name"`
	UIPort  int       `yaml:"ui_port"`
	GQLPort int       `yaml:"gql_port"`
	Author  string    `yaml:"author"`
	Date    string    `yaml:"date"`
	Modules []*Module `yaml:"modules"`
}
