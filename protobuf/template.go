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
{{ range .Types}}
message {{pascal .Name}} {
	int64 id = 1;
    {{- range $index, $elem := .Fields}}
    {{typeTrans "protobuf" $elem}} {{$elem.Name}} = {{add $index 2}};
    {{- end}}
}

message Get{{pascal .Name}}Input {
	string id = 1;
}

message {{pascal .Name}}Input {
    {{- range $index, $elem := .Fields}}
    {{typeTrans "protobuf" $elem}} {{$elem.Name}} = {{add $index 1}};
    {{- end}}
}

message {{pluralize (pascal .Name)}} {
	repeated {{pascal .Name}} {{pluralize (lower .Name)}} = 1;
}


{{- end}}

service {{pascal .Name}} {
    {{- range .Funcs}}
	rpc {{pascal .Name}} ({{pascal .Name}}Input) returns ({{pascal .Name}}Payload); 
    {{- end}}
	{{- range .Types}}
	rpc Get (Get{{pascal .Name}}Input) returns ({{pascal .Name}}); 
	rpc List (Empty) returns ({{pluralize (pascal .Name)}}); 
	{{- end}}
}

{{- range .Funcs}}
message {{pascal .Name}}Input {
	{{- if .Params}}
    {{- range $index, $elem := .Params}}
    {{inputTrans "protobuf" $elem}} {{$elem.Name}} = {{add $index 1}};
    {{- end}}
	{{- end}}
}

message {{pascal .Name}}Payload {
	{{- if .Return}}
    {{- range $index, $elem := .Return}}
    {{typeTrans "protobuf" $elem}} {{$elem.Name}} = {{add $index 1}};
	{{- end}}
	{{- end}}
}
{{- end}}
`)
