package core

type System struct {
	Name    string
	Modules []*Module
}

func (s *System) Walk(v Visitor) error {
	for _, m := range s.Modules {
		err := v.Visit(m)
		if err != nil {
			return err
		}
	}
	return nil
}
