package lang

import (
	"github.com/jt05610/scaf/system"
	"os"
	"path/filepath"
	"text/template"
)

type File struct {
	Path    string
	TmpFunc func() *template.Template
}

func Generate(parDir string, ff func(string) []*File, m *system.Module) error {
	err := os.MkdirAll(filepath.Join(parDir, "src"), 0755)
	if err != nil {
		return err
	}
	for _, f := range ff(parDir) {
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
