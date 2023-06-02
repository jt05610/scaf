package ts

import (
	"github.com/jt05610/scaf/system"
	"html/template"
	"os"
	"path/filepath"
)

type generator struct {
	parDir string
}

func (g *generator) pkgJson(m system.Module) error {
	f, err := os.Create(filepath.Join(g.parDir, "package.json"))
	if err != nil {
		return err
	}
	pkg := packageJson()
	return pkg.Execute(f, m)
}

type file struct {
	Path    string
	TmpFunc func() *template.Template
}

func (g *generator) Generate(m *system.Module) error {
	for _, f := range []file{
		{
			Path:    filepath.Join(g.parDir, "package.json"),
			TmpFunc: packageJson,
		},
	} {
		fp, err := os.Create(f.Path)
		if err != nil {
			return err
		}
		err = f.TmpFunc().Execute(fp, m)
		if err != nil {
			return err
		}
		err = fp.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func NewGenerator(parDir string) system.ModuleGenerator {
	return &generator{parDir: parDir}
}
