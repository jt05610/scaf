package core

import (
	"embed"
	"fmt"
)

type Scripts struct {
	Init  string
	Gen   string
	Start string
	Stop  string
}

type TypeMap struct {
	Int    string
	Float  string
	String string
	Bool   string
}

func (m *TypeMap) Map(t BaseType) (string, bool) {
	switch t {
	case IntType:
		return m.Int, true
	case FloatType:
		return m.Float, true
	case StringType:
		return m.String, true
	case BoolType:
		return m.Bool, true
	default:
		return "", false
	}
}

type Language struct {
	*Cmd     `yaml:"-"`
	Service  string    `yaml:"service"`
	Name     string    `yaml:"name"`
	Scripts  *Scripts  `yaml:"scripts"`
	fs       *embed.FS `yaml:"fs"`
	TypeMap  *TypeMap  `yaml:"type_map"`
	ArrayFmt string    `yaml:"array_fmt"`
}

func (l *Language) FS() *embed.FS {
	return l.fs
}

func (l *Language) MapType(t BaseType) (string, bool) {
	s, b := l.TypeMap.Map(t)
	return s, b
}

func (l *Language) MakeArray(s string) string {
	return fmt.Sprintf(l.ArrayFmt, s)
}

func CreateLanguage(name, parent string, scripts *Scripts, fs *embed.FS, types *TypeMap, arrayFmt string) *Language {
	return &Language{
		Name:     name,
		Cmd:      NewCmd(parent, scripts),
		fs:       fs,
		TypeMap:  types,
		ArrayFmt: arrayFmt,
	}
}
