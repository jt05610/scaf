package solidjs

type Host struct {
	Remotes []*Remote
	Shared  []string
	Port    int
}

func (h *Host) TplDir() string {
	return "host/template"
}
