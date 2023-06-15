package core

type System struct {
	Name    string
	GQLPort int
	Author  string
	Date    string
	Modules []*Module
}

func (s *System) Walk(v Visitor) error {
	for _, m := range s.Modules {
		err := v.VisitModule(m)
		if err != nil {
			return err
		}
	}
	return nil
}
