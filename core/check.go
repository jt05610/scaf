package core

type Checker struct {
	seen map[string]int
	k    int
}

func (c *Checker) VisitAPI(a *API) error {
	for _, d := range a.Deps {
		if err := c.VisitAPI(d); err != nil {
			return err
		}
	}
	return nil
}

func (c *Checker) VisitModule(m *Module) error {
	if _, seen := c.seen[m.Name]; seen {
		return nil
	}
	c.k++
	c.seen[m.Name] = c.k

	return nil
}

var i = `
func (c *Checker) IsAcyclic(s *System) bool {
	c.seen = make(map[string]int)
	_ = s.Walk(c)
	for _, mm := range s.Modules {
		po := c.seen[mm.Name]
		for _, m := range mm.Deps {
			if po <= c.seen[m.Name] {
				return false
			}
		}
	}
	return true
}
`

func (c *Checker) PostOrder(name string) (int, bool) {
	o, seen := c.seen[name]
	return o, seen
}

func NewChecker() *Checker {
	return &Checker{seen: make(map[string]int)}
}
