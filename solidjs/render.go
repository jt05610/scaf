package solidjs

import (
	"errors"
	"github.com/jt05610/scaf/core"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const MaxRecurseLevel = 10

type hostRenderer struct {
	outDir string
}

func (r *hostRenderer) Visit(m *core.Module) core.Visitor {
	//TODO implement me
	panic("implement me")
}

func renderFile(name, outDir string, h *Host) error {
	tmpl, err := template.ParseFiles(name)
	if err != nil {
		return err
	}
	outName := strings.TrimSuffix(name, ".gotpl")
	outName = strings.Replace(outName, h.TplDir(), outDir, 1)
	f, err := os.Create(outName)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			panic(err)
		}
	}(f)
	return tmpl.Execute(f, h)
}

func ApplyToDir(dir string, action func(string) error) {

}
func (r *hostRenderer) renderDir(level int, name string, h *Host) error {
	if level > MaxRecurseLevel {
		return errors.New("max recursion level reached")
	}
	files, err := os.ReadDir(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(strings.Replace(name, h.TplDir(), r.outDir, 1), os.ModePerm)
	if err != nil {
		return err
	}
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".gotpl" {
			err := renderFile(name+"/"+file.Name(), r.outDir, h)
			if err != nil {
				return err
			}
		} else if file.IsDir() && level < MaxRecurseLevel {
			err := r.renderDir(level+1, name+"/"+file.Name(), h)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *hostRenderer) Render(_ io.Writer, h *Host) error {
	return r.renderDir(0, h.TplDir(), h)
}

func NewHostRenderer(outDir string) core.Visitor {
	return &hostRenderer{outDir: outDir}
}
