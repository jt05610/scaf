package solidjs

type PublicAsset struct {
	Name string
	Path string
}

type Remote struct {
	Name    string
	Addr    string
	Port    int
	Exposes []*PublicAsset
	Shared  []string
}

func (r *Remote) TplDir() string {
	return "remote/tpl"
}
