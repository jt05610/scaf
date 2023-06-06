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

func (d *dirBuilder) buildDir(parent string, m *core.Module, dir *Dir) {
	d.mu.Lock()
	if _, seen := d.seen[dir]; seen {
		return
	}
	d.seen[dir] = true
	d.mu.Unlock()
	path := filepath.Join(parent, dir.Name(m))
	err := os.MkdirAll(path, 0755)
	if err != nil {
		panic(err)
	}
	var wg sync.WaitGroup
	for _, file := range dir.Files {
		wg.Add(1)
		go func(file *File) {
			defer wg.Done()
			fp := filepath.Join(path, file.Name(m))
			f, err := os.Create(fp)
			if err != nil {
				panic(err)
			}
			defer func() {
				_ = f.Close()
			}()
			if err := file.Template.Execute(f, m); err != nil {
				panic(err)
			}
		}(file)
	}
	for _, child := range dir.Children {
		wg.Add(1)
		go func(child *Dir) {
			defer wg.Done()
			d.buildDir(path, m, child)
		}(child)
	}
	wg.Wait()
}

func (d *dirBuilder) Visit(m *core.Module) core.Visitor {
	path := filepath.Join(d.parent, m.Name)
	var wg sync.WaitGroup
	for _, dir := range d.dirs {
		wg.Add(1)
		go func(dir *Dir) {
			defer wg.Done()
			d.buildDir(path, m, dir)
		}(dir)
	}
	wg.Wait()
	return d
}

func NewDirBuilder(dirs ...*Dir) core.Visitor {
	return &dirBuilder{dirs: dirs, seen: make(map[*Dir]bool)}
}
