package protobuf

import (
	"github.com/jt05610/scaf/builder"
)

var protoTemplate = builder.NewTemplate(
	`// Auto-generated protobuf schema. Do not edit. Generated on {{.Date}} by github.com/jt05610/scaf v0.1.0

syntax = "proto3";

package {{.Name}};

option go_package = "{{.Name}}/v1/{{.Name}}";

message Empty {}

{{- range .Types}}
message {{pascal .Name}} {
    {{- range $index, $elem := .Fields}}
    {{typeTrans "protobuf" $elem}} {{$elem.Name}} = {{add $index 1}};
    {{- end}}
}

message {{pluralize (pascal .Name)}} {
	repeated {{pascal .Name}} {{pluralize (lower .Name)}} = 1;
}
{{- end}}

service {{pascal .Name}}Service {
    {{- range .Funcs}}
    rpc {{pascal .Name}} ({{pascal .Name}}Request) returns ({{pascal .Name}}Response) {}
    {{- end}}
}

{{- range .Funcs}}
message {{pascal .Name}}Request {
	{{- if .Params}}
    {{- range $index, $elem := .Params}}
    {{typeTrans "protobuf" $elem}} {{$elem.Name}} = {{add $index 1}};
    {{- end}}
	{{- end}}
}

message {{pascal .Name}}Response {
	{{- if .Return}}
    {{- range $index, $elem := .Return}}
    {{typeTrans "protobuf" $elem}} {{$elem.Name}} = {{add $index 1}};
	{{- end}}
	{{- end}}
}
{{- end}}
`)
