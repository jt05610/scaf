package core

type BaseType uint

const (
	Int BaseType = iota
	Float
	String
	Bool
	Custom
)

type Language string

const (
	GQL        Language = "gql"
	Protobuf   Language = "protobuf"
	Go         Language = "go"
	Python     Language = "python"
	Fortran    Language = "fortran"
	C          Language = "c"
	Cpp        Language = "cpp"
	Typescript Language = "typescript"
	Clojure    Language = "clojure"
)

var TypeMapping = map[Language][]string{
	GQL: {
		Int:    "Int",
		Float:  "Float",
		String: "String",
		Bool:   "Boolean",
	},
	Protobuf: {
		Int:    "int32",
		Float:  "float",
		String: "string",
		Bool:   "bool",
	},
	Go: {
		Int:    "int",
		Float:  "float32",
		String: "string",
		Bool:   "bool",
	},
	Python: {
		Int:    "int",
		Float:  "float",
		String: "str",
		Bool:   "bool",
	},
	Fortran: {
		Int:    "INTEGER",
		Float:  "REAL",
		String: "CHARACTER",
		Bool:   "LOGICAL",
	},
	C: {
		Int:    "int",
		Float:  "float",
		String: "char*",
		Bool:   "bool",
	},
	Cpp: {
		Int:    "int",
		Float:  "float",
		String: "std::string",
		Bool:   "bool",
	},
	Typescript: {
		Int:    "number",
		Float:  "number",
		String: "string",
		Bool:   "boolean",
	},
	Clojure: {
		Int:    "integer",
		Float:  "float",
		String: "string",
		Bool:   "boolean",
	},
}

type Field struct {
	Name       string
	Type       BaseType
	CustomType string
	Last       bool
}

func (f *Field) TypeString(l Language) string {
	if f.Type == Custom {
		return f.CustomType
	}
	return TypeMapping[l][f.Type]
}

type FieldList []*Field

type Type struct {
	Name      string
	HasPlural bool
	Fields    FieldList
	Funcs     []*Func
}

type Func struct {
	Name   string
	Return FieldList
	Params FieldList
}

type Module struct {
	Name     string
	Date     string
	Port     int
	Language Language
	Deps     []*Module
	Types    []*Type
	Funcs    []*Func
}
