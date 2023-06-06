package core

type System struct {
	Modules []*Module
}

func (s *System) Walk(v Visitor) Visitor {
	for _, m := range s.Modules {
		v.Visit(m)
	}
	return v
}
