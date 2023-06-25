package lang

import (
	"embed"
	"fmt"
	"github.com/jt05610/scaf/core"
)

type Scripts struct {
	WorkDir string            `yaml:"work_dir"`
	Map     map[string]string `yaml:"inline"`
}

type TypeMap map[core.BaseType]string

type Language struct {
	*CmdSet  `yaml:"-"`
	Service  string    `yaml:"service"`
	Name     string    `yaml:"name"`
	Scripts  *Scripts  `yaml:"scripts"`
	fs       *embed.FS `yaml:"fs"`
	TypeMap  TypeMap   `yaml:"type_map"`
	ArrayFmt string    `yaml:"array_fmt"`
}

func (l *Language) FS() *embed.FS {
	return l.fs
}

func (l *Language) TypeString(t core.Type) (s string) {
	if t.IsPrimitive() {
		s = l.TypeMap[core.BaseType(t.String())]
	} else {
		s = t.String()
	}
	return s
}

func (l *Language) InputString(t core.Type) (s string) {
	if t.IsPrimitive() {
		s = l.TypeMap[core.BaseType(t.String())]
	} else {
		s = t.String() + "Params"
	}
	return s
}

func (l *Language) MakeArray(s string) string {
	return fmt.Sprintf(l.ArrayFmt, s)
}

func (l *Language) TypeDecl(f core.Field) (s string) {
	s = l.TypeString(f.Type)
	if f.IsArray {
		s = l.MakeArray(s)
	}
	return s
}

func (l *Language) InputDecl(f core.Field) (s string) {
	s = l.InputString(f.Type)
	if f.IsArray {
		s = l.MakeArray(s)
	}
	return s
}

func (l *Language) CreateDecl(f core.Field) (s string) {
	s = l.InputString(f.Type)
	if !f.Type.IsPrimitive() {
		s = "Create" + s
	}
	if f.IsArray {
		s = l.MakeArray(s)
	}

	return s
}
func (l *Language) UpdateDecl(f core.Field) (s string) {
	s = l.InputString(f.Type)
	if !f.Type.IsPrimitive() {
		s = "Update" + s
	}
	if f.IsArray {
		s = l.MakeArray(s)
	}

	return s
}
func CreateLanguage(name, parent string, sysScripts, modScripts *Scripts, fs *embed.FS, types TypeMap, arrayFmt string) *Language {
	return &Language{
		Name:     name,
		CmdSet:   NewCmdSet(parent, sysScripts, modScripts),
		Scripts:  sysScripts,
		fs:       fs,
		TypeMap:  types,
		ArrayFmt: arrayFmt,
	}
}
