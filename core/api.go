package core

type Lang string

const (
	GQL      Lang = "gql"
	TS       Lang = "ts"
	Protobuf Lang = "proto"
	Go       Lang = "go"
	Python   Lang = "py"
)

type API struct {
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	PortMap     *PortMap `yaml:"-"`
	Author      string   `yaml:"author"`
	Version     int      `yaml:"version"`
	Language    Lang     `yaml:"language"`
	Date        string   `yaml:"date"`
	Deps        []*API   `yaml:"deps"`
	Models      []*Model `yaml:"types"`
	Funcs       []*Func  `yaml:"funcs"`
}

func (a *API) HasSubs() bool {
	for _, m := range a.Models {
		for _, f := range m.Fields {
			if f.Subscribe {
				return true
			}
		}
	}
	return false
}

func NewAPI(name, author, date string, lang Lang) *API {
	return &API{
		Name:     name,
		Author:   author,
		Language: lang,
		Date:     date,
		Deps:     make([]*API, 0),
		Models:   make([]*Model, 0),
		Funcs:    make([]*Func, 0),
	}
}

func (a *API) AddModel(t *Model) {
	a.Models = append(a.Models, t)
}

func (a *API) AddFunc(f *Func) {
	a.Funcs = append(a.Funcs, f)
}

func (a *API) AddDep(d *API) {
	a.Deps = append(a.Deps, d)
}
