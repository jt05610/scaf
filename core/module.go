package core

import "fmt"

type BaseType string

const (
	IntType    BaseType = "int"
	FloatType  BaseType = "float"
	StringType BaseType = "string"
	BoolType   BaseType = "bool"
)

type Language string

const (
	GQL        Language = "gql"
	Protobuf   Language = "pb"
	Go         Language = "go"
	Python     Language = "python"
	Fortran    Language = "fortran"
	C          Language = "c"
	Cpp        Language = "cpp"
	Typescript Language = "typescript"
	Clojure    Language = "clojure"
	Rust       Language = "rust"
)

type Type struct {
	Name      string
	Fields    FieldList
	Funcs     []*Func
	isArray   bool
	Query     bool
	Subscribe bool
}

var Int = &Type{Name: "int"}
var Float = &Type{Name: "float"}
var String = &Type{Name: "string"}
var Bool = &Type{Name: "bool"}

func Array(t *Type) *Type {
	t.isArray = true
	return t
}

func makeArray(l Language, s string) string {
	switch l {
	case GQL:
		return fmt.Sprintf("[%s]", s)
	case Protobuf:
		return fmt.Sprintf("repeated %s", s)
	case Go, Clojure:
		return fmt.Sprintf("[]%s", s)
	case Python:
		return fmt.Sprintf("List[%s]", s)
	case Fortran:
		return fmt.Sprintf("type(%s), allocatable, dimension(:)", s)
	case C:
		return fmt.Sprintf("%s *", s)
	case Cpp:
		return fmt.Sprintf("std::vector<%s>", s)
	case Typescript:
		return fmt.Sprintf("%s[]", s)
	case Rust:
		return fmt.Sprintf("Vec<%s>", s)
	default:
		return s
	}
}

func (t *Type) string(l Language, isInput bool) (ret string) {
	if _, found := TypeMapping[l][BaseType(t.Name)]; found {
		ret = TypeMapping[l][BaseType(t.Name)]
	} else {
		ret = t.Name
		if isInput {
			ret += "Input"
		}
	}
	if t.isArray {
		return makeArray(l, ret)
	}
	return

}

func (t *Type) String(l Language) string {
	return t.string(l, false)
}

func (t *Type) InputString(l Language) string {
	return t.string(l, true)
}

var TypeMapping = map[Language]map[BaseType]string{
	GQL: {
		IntType:    "Int",
		FloatType:  "Float",
		StringType: "String",
		BoolType:   "Boolean",
	},
	Protobuf: {
		IntType:    "int32",
		FloatType:  "float",
		StringType: "string",
		BoolType:   "bool",
	},
	Go: {
		IntType:    "int",
		FloatType:  "float32",
		StringType: "string",
		BoolType:   "bool",
	},
	Python: {
		IntType:    "int",
		FloatType:  "float",
		StringType: "str",
		BoolType:   "bool",
	},
	Fortran: {
		IntType:    "INTEGER",
		FloatType:  "REAL",
		StringType: "CHARACTER",
		BoolType:   "LOGICAL",
	},
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
	Typescript: {
		IntType:    "number",
		FloatType:  "number",
		StringType: "string",
		BoolType:   "boolean",
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
