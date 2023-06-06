package core

type BaseType string

const (
	Int    BaseType = "Int"
	Float  BaseType = "Float"
	String BaseType = "String"
	Bool   BaseType = "Bool"
)

type Field struct {
	Name string
	Type BaseType
	Last bool
}

type FieldList []*Field

type Type struct {
	Name   string
	Fields FieldList
	Funcs  []*Func
}

type Func struct {
	Name   string
	Return FieldList
	Params FieldList
}

type Module struct {
	Name  string
	Date  string
	Port  int
	Deps  []*Module
	Types []*Type
	Funcs []*Func
}
