package gql

import (
	"github.com/jt05610/scaf/builder"
	"text/template"
)

var toolsTemplate = template.Must(
	template.New("tools").
		Parse(`//go:build tools
// +build tools

package tools

import (
	_ "github.com/99designs/gqlgen"
	_ "github.com/99designs/gqlgen/graphql/introspection"
)
`))

var gqlTemplate = builder.NewTemplate(`# Auto-generated GraphQL schema. Do not edit. Generated on {{.Date}} by github.com/jt05610/scaf v0.1.0

schema {
    query: Query
    mutation: Mutation
    # subscription: Subscription
}

{{- range .Types}}
type {{.Name}} {
	id : ID!
    {{- range .Fields}}
    {{.Name}}: {{typeTrans "gql" .}}
    {{- end}}
}

type {{pluralize (pascal .Name)}} {
	{{lower (pluralize .Name)}}: [{{.Name}}]
}

input {{.Name}}Input {
    {{- range .Fields}}
    {{.Name}}: {{typeTrans "gql" .}}
    {{- end}}
}

{{- end}}

type Query {
    {{- range .Types}}
{{- if .Query}}
    {{lower .Name}}(id: ID!): {{.Name}}
    {{pluralize (lower .Name)}}: [{{.Name}}]
{{- end}}
    {{- end}}
}
{{- range .Funcs}}

input {{.Name}}Input {
	{{- range .Params}}
	{{.Name}}: {{inputTrans "gql" .}}
	{{- end}}	
}

type {{.Name}}Payload {
	{{- range .Return}}
	{{.Name}}: {{typeTrans "gql" .}}
	{{- end}}
}

{{- end}}

type Mutation {
    {{- range .Funcs}}
    {{.Name}}(input: {{.Name}}Input): {{.Name}}Payload
    {{- end}}
}
`)

var resolverTemplate = builder.NewTemplate(`// Auto-generated GraphQL resolver. Do not edit. Generated on {{.Date}} by github.com/jt05610/scaf v0.1.0
package graph

import "{{.Name}}/v1"

type Resolver struct{
	client {{.Name}}.{{pascal .Name}}Client
}
`)

var schemaResolverTemplate = builder.NewTemplate(`// Auto-generated GraphQL schema resolver. Do not edit. Generated on {{.Date}} by github.com/jt05610/scaf v0.1.0
package graph
{{- $pkgName := .Name}}
import (
	"context"
	"fmt"
	"{{.Name}}/v1"
)

{{- range .Funcs}}
// {{.Name}} is the resolver for the {{.Name}} field.
func (r *mutationResolver) {{.Name}}(ctx context.Context, input *{{$pkgName}}.{{.Name}}Input) (*housework.{{.Name}}Payload, error) {
	return r.client.{{.Name}}(ctx, input)
}

{{- end}}

{{- range .Types}}
{{- if .Query}}
// {{.Name}} is the resolver for the {{.Name}} field.
func (r *queryResolver) {{.Name}}(ctx context.Context, id string) (*{{$pkgName}}.{{.Name}}, error) {
	return r.client.Get(ctx, &{{$pkgName}}.Get{{.Name}}Input{Id: id})
}

// {{pluralize .Name}} is the resolver for the {{pluralize .Name}} field.
func (r *queryResolver) {{pluralize .Name}}(ctx context.Context) ([]*{{$pkgName}}.{{.Name}}, error) {
	c, err := r.client.List(ctx, nil)
	if err != nil {
		return nil, err
	}
	return c.{{pluralize .Name}}, nil
}

{{- end}}
{{- end}}
`)

var serverTemplateInstance = builder.NewTemplate(`package main

import (
	"log"
	"net/http"
	"{{.Name}}/graph"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

func main() {
	
	port := "{{.Port}}"

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
`)

var gqlGenYmlTemplate = builder.NewTemplate(`schema:
  - graph/*.graphql

exec:
  filename: graph/generated.go
  package: graph

federation:
  filename: graph/federation.go
  package: graph

resolver:
  layout: follow-schema
  dir: graph
  package: graph
  filename_template: "{name}.resolvers.go"

autobind:
  - "{{.Name}}/v1"

models:
  ID:
    model:
      - github.com/99designs/gqlgen/graphql.ID
      - github.com/99designs/gqlgen/graphql.Int
      - github.com/99designs/gqlgen/graphql.Int64
      - github.com/99designs/gqlgen/graphql.Int32
  Int:
    model:
      - github.com/99designs/gqlgen/graphql.Int
      - github.com/99designs/gqlgen/graphql.Int64
      - github.com/99designs/gqlgen/graphql.Int32

`)
