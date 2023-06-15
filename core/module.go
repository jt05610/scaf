package core

type BaseType string

const (
	IntType    BaseType = "int"
	FloatType  BaseType = "float"
	StringType BaseType = "string"
	BoolType   BaseType = "bool"
)

type Type struct {
	Name      string
	Plural    string
	Fields    FieldList
	Funcs     []*Func
	isArray   bool
	Query     bool
	Subscribe bool
}

type Aggregate struct {
	Type
	Individual *Type
}

var Int = &Type{Name: "int"}
var Float = &Type{Name: "float"}
var String = &Type{Name: "string"}
var Bool = &Type{Name: "bool"}

func Array(t *Type) *Type {
	t.isArray = true
	return t
}
func (t *Type) string(l Language, isInput bool) (ret string) {
	var found bool
	if ret, found = l.MapType(BaseType(t.Name)); !found {
		ret = t.Name
		if isInput {
			ret += "Input"
		}
	}
	if t.isArray {
		return l.MakeArray(ret)
	}
	return ret
}

func (t *Type) String(l Language) string {
	return t.string(l, false)
}

func (t *Type) InputString(l Language) string {
	return t.string(l, true)
}

var old = `

func makeArray(l Language, s string) string {
	switch l {
	case Clojure:
		return fmt.Sprintf("[]%s", s)
	case C:
		return fmt.Sprintf("%s *", s)
	case Cpp:
		return fmt.Sprintf("std::vector<%s>", s)
	case Rust:
		return fmt.Sprintf("Vec<%s>", s)
	default:
		return s
	}
}

var TypeMapping = map[Language]map[BaseType]string{
	C: {
		IntType:    "int",
		FloatType:  "float",
		StringType: "char*",
		BoolType:   "bool",
	},
	Cpp: {
		IntType:    "int",
		FloatType:  "float",
		StringType: "std::string",
		BoolType:   "bool",
	},
	Clojure: {
		IntType:    "integer",
		FloatType:  "float",
		StringType: "string",
		BoolType:   "boolean",
	},
	Rust: {
		IntType:    "i32",
		FloatType:  "f32",
		StringType: "String",
		BoolType:   "bool",
	},
}
`

type Field struct {
	Name string
	Type *Type
	Last bool
}

func (f *Field) TypeString(l Language) string {
	return f.Type.String(l)
}

func (f *Field) InputString(l Language) string {
	return f.Type.InputString(l)
}

type FieldList []*Field
type Func struct {
	Name   string
	Return FieldList
	Params FieldList
}

type Module struct {
	Name     string
	Author   string
	Version  int
	Date     string
	Port     int
	Language Language
	Deps     []*Module
	Types    []*Type
	Funcs    []*Func
}
