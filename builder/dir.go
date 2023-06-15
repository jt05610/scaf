package builder

import (
	"github.com/jt05610/scaf/core"
	"os"
	"path/filepath"
	"sync"
	"text/template"
)

type File struct {
	Name     func(*core.Module) string
	Template *template.Template
}

type Dir struct {
	Name     func(*core.Module) string
	Children []*Dir
	Files    []*File
}

type dirBuilder struct {
	parent string
	dirs   []*Dir
	mu     sync.RWMutex
	seen   map[*Dir]bool
}

func NewDir(nameFunc func(*core.Module) string, files []*File) *Dir {
	return &Dir{Name: nameFunc, Files: files}
}

func (d *Dir) AddChild(child *Dir) {
	d.Children = append(d.Children, child)
}

func (d *dirBuilder) buildDir(parent string, m *core.Module, dir *Dir) error {
	d.mu.Lock()
	if _, seen := d.seen[dir]; seen {
		return nil
	}
	d.seen[dir] = true
	d.mu.Unlock()
	path := filepath.Join(parent, dir.Name(m))
	err := os.MkdirAll(path, 0755)
	if err != nil {
		return err
	}
	for _, file := range dir.Files {
		fp := filepath.Join(path, file.Name(m))
		f, err := os.Create(fp)
		if err != nil {
			return err
		}
		if err := file.Template.Execute(f, m); err != nil {
			_ = f.Close()
			return err
		}
		_ = f.Close()
	}
	for _, child := range dir.Children {
		if err := d.buildDir(path, m, child); err != nil {
			return err
		}
	}
	return nil
}

func (d *dirBuilder) VisitModule(m *core.Module) error {
	path := filepath.Join(d.parent, m.Name)
	for _, dir := range d.dirs {
		if err := d.buildDir(path, m, dir); err != nil {
			return err
		}
	}
	return nil
}
