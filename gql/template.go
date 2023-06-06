package gql

import (
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

var gqlTemplate = template.Must(
	template.New("gql").
		Parse(`# Auto-generated GraphQL schema. Do not edit. Generated on {{.Date}} by github.com/jt05610/scaf v0.1.0

schema {
    query: Query
    mutation: Mutation
    subscription: Subscription
}

{{- range .Types}}
type {{.Name}} {
    {{- range .Fields}}
    {{.Name}}: {{.Type}}
    {{- end}}
}
{{- end}}

type Query {
    {{- range .Types}}
    get{{.Name}}(id: ID!): {{.Name}}
    list{{.Name}}s: [{{.Name}}]
    {{- end}}
}

type Mutation {
    {{- range .Types}}
    create{{.Name}}({{range .Fields}}{{.Name}}: {{.Type}}{{if not .Last}},{{end}}{{end}}): {{.Name}}
    update{{.Name}}(id: ID!, {{range .Fields}}{{.Name}}: {{.Type}}{{if not .Last}},{{end}}{{end}}): {{.Name}}
    delete{{.Name}}(id: ID!): {{.Name}}
    {{- end}}
    {{- range .Funcs}}
    {{.Name}}({{range .Params}}{{.Name}}: {{.Type}}{{if not .Last}},{{end}}{{end}}): {{range .Return}}{{.Type}}{{if not .Last}},{{end}}{{end}}
    {{- end}}
}

type Subscription {
    {{- range .Types}}
    on{{.Name}}Created: {{.Name}}
    on{{.Name}}Updated: {{.Name}}
    on{{.Name}}Deleted: {{.Name}}
    {{- end}}
}
`))

var resolverTemplate = template.Must(
	template.New("resolver").
		Parse(`// Auto-generated GraphQL resolver. Do not edit. Generated on {{.Date}} by github.com/jt05610/scaf v0.1.0
package {{.Name}}

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct{
	{{- range .Types}}
	{{.Name}}s []*{{.Name}}
	{{- end}}
}
`))

var schemaResolverTemplate = template.Must(
	template.New("schemaResolver").
		Parse(`// Auto-generated GraphQL schema resolver. Do not edit. Generated on {{.Date}} by github.com/jt05610/scaf v0.1.0
package {{.Name}}

import (
	"context"
	"fmt"
	"github.com/jt05610/scaf/core"
	"github.com/jt05610/scaf/gql"
	"github.com/jt05610/scaf/gql/generated"
	"github.com/jt05610/scaf/gql/model"
	"github.com/jt05610/scaf/gql/resolver"
	"github.com/jt05610/scaf/gql/subscription"
	"github.com/jt05610/scaf/gql/validator"
	"github.com/jt05610/scaf/log"
	"github.com/jt05610/scaf/store"
	"github.com/jt05610/scaf/util"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"github.com/vektah/gqlparser/v2/parser"
	"github.com/vektah/gqlparser/v2/validator"
	"io/ioutil"	
	"os"
	"path/filepath"
	"strings"
	"time"
)	

// This file will not be regenerated automatically.
//
	
func NewSchemaResolver() *SchemaResolver {
	return &SchemaResolver{
		Resolver: resolver.Resolver{
			{{- range .Types}}
			{{.Name}}s: []*model.{{.Name}}{},
			{{- end}}
		},
	}
}

type SchemaResolver struct {
	resolver.Resolver
}

func (r *SchemaResolver) Query() generated.QueryResolver {
	return &queryResolver{r}
}

func (r *SchemaResolver) Mutation() generated.MutationResolver {
	return &mutationResolver{r}
}
	
func (r *SchemaResolver) Subscription() generated.SubscriptionResolver {
	return &subscriptionResolver{r}
}

type queryResolver struct{ *SchemaResolver }

{{- range .Types}}
func (r *queryResolver) Get{{.Name}}(ctx context.Context, id string) (*model.{{.Name}}, error) {
	{{.Name}}, err := r.Get{{.Name}}ByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return {{.Name}}, nil
}

func (r *queryResolver) List{{.Name}}s(ctx context.Context) ([]*model.{{.Name}}, error) {
	return r.List{{.Name}}s(ctx)
}
{{- end}}

type mutationResolver struct{ *SchemaResolver }

{{- range .Types}}

func (r *mutationResolver) Create{{.Name}}(ctx context.Context, {{range .Fields}}{{.Name}} {{.Type}}{{if not .Last}}, {{end}}{{end}}) (*model.{{.Name}}, error) {
	{{.Name}}, err := r.Create{{.Name}}(ctx, {{range .Fields}}{{.Name}}{{if not .Last}}, {{end}}{{end}})
	if err != nil {
		return nil, err
	}
	return {{.Name}}, nil
}

func (r *mutationResolver) Update{{.Name}}(ctx context.Context, id string, {{range .Fields}}{{.Name}} {{.Type}}{{if not .Last}}, {{end}}{{end}}) (*model.{{.Name}}, error) {
	{{.Name}}, err := r.Update{{.Name}}(ctx, id, {{range .Fields}}{{.Name}}{{if not .Last}}, {{end}}{{end}})
	if err != nil {
		return nil, err
	}
	return {{.Name}}, nil
}

func (r *mutationResolver) Delete{{.Name}}(ctx context.Context, id string) (*model.{{.Name}}, error) {
	{{.Name}}, err := r.Delete{{.Name}}(ctx, id)
	if err != nil {
		return nil, err
	}
	return {{.Name}}, nil
}
{{- end}}

type subscriptionResolver struct{ *SchemaResolver }

{{- range .Types}}

func (r *subscriptionResolver) On{{.Name}}Created(ctx context.Context) (<-chan *model.{{.Name}}, error) {
	return subscription.Subscribe{{.Name}}Created(ctx)
}	

func (r *subscriptionResolver) On{{.Name}}Updated(ctx context.Context) (<-chan *model.{{.Name}}, error) {
	return subscription.Subscribe{{.Name}}Updated(ctx)
}

func (r *subscriptionResolver) On{{.Name}}Deleted(ctx context.Context) (<-chan *model.{{.Name}}, error) {
	return subscription.Subscribe{{.Name}}Deleted(ctx)
}
{{- end}}

func (r *SchemaResolver) LoadSchema(ctx context.Context) error {
	// Load schema
	schema, err := r.LoadSchemaFromFile(ctx)
	if err != nil {
		return err
	}
	
	// Validate schema
	err = r.ValidateSchema(ctx, schema)
	if err != nil {
		return err
	}
	
	// Load resolvers
	err = r.LoadResolvers(ctx, schema)
	if err != nil {
		return err
	}
	
	// Load subscriptions
	err = r.LoadSubscriptions(ctx, schema)
	if err != nil {
		return err
	}
	
	return nil
}

func (r *SchemaResolver) LoadSchemaFromFile(ctx context.Context) (*gql.Schema, error) {
	// Load schema file
	schemaFile, err := os.Open("schema.graphql")
	if err != nil {
		return nil, err
	}
	defer schemaFile.Close()
	
	// Read schema file
	schemaBytes, err := ioutil.ReadAll(schemaFile)
	if err != nil {
		return nil, err
	}
	
	// Parse schema
	schema, err := parser.ParseSchema(&ast.Source{Name: "schema.graphql", Input: string(schemaBytes)})
	if err != nil {
		return nil, err
	}
	
	return schema, nil
}

func (r *SchemaResolver) ValidateSchema(ctx context.Context, schema *gql.Schema) error {
	// Validate schema
	err := validator.ValidateSchema(schema)
	if err != nil {
		return err
	}
	
	return nil
}

func (r *SchemaResolver) LoadResolvers(ctx context.Context, schema *gql.Schema) error {
	// Load resolvers
	for _, schemaType := range schema.Types {
		// Skip system types
		if schemaType.BuiltIn {
			continue
		}
		
		// Load resolver
		err := r.LoadResolver(ctx, schemaType)
		if err != nil {
			return err
		}
	}
	
	return nil
}

func (r *SchemaResolver) LoadResolver(ctx context.Context, schemaType *gql.Type) error {
	// Load resolver
	resolver, err := r.LoadResolverFromFile(ctx, schemaType)
	if err != nil {
		return err
	}
	
	// Set resolver
	r.SetResolver(schemaType.Name, resolver)
	
	return nil
}

func (r *SchemaResolver) LoadResolverFromFile(ctx context.Context, schemaType *gql.Type) (interface{}, error) {
	// Load resolver file
	resolverFile, err := os.Open(filepath.Join("resolvers", strings.ToLower(schemaType.Name) + ".go"))
	if err != nil {
		return nil, err
	}
	defer resolverFile.Close()
	
	// Read resolver file
	resolverBytes, err := ioutil.ReadAll(resolverFile)
	if err != nil {
		return nil, err
	}
	
	// Compile resolver
	resolver, err := r.CompileResolver(ctx, schemaType, string(resolverBytes))
	if err != nil {
		return nil, err
	}
	
	return resolver, nil
}

func (r *SchemaResolver) CompileResolver(ctx context.Context, schemaType *gql.Type, resolverCode string) (interface{}, error) {
	// Create resolver context
	resolverContext := &ResolverContext{
		Schema: schemaType,
	}
	
	// Create resolver
	resolver, err := r.CreateResolver(ctx, resolverContext, resolverCode)
	if err != nil {
		return nil, err
	}
	
	return resolver, nil
}

func (r *SchemaResolver) CreateResolver(ctx context.Context, resolverContext *ResolverContext, resolverCode string) (interface{}, error) {
	// Create resolver
	resolver, err := vm.CreateResolver(resolverContext, resolverCode)
	if err != nil {
		return nil, err
	}
	
	return resolver, nil
}

func (r *SchemaResolver) SetResolver(name string, resolver interface{}) {
	// Set resolver
	r.Resolvers[name] = resolver
}

func (r *SchemaResolver) LoadSubscriptions(ctx context.Context, schema *gql.Schema) error {
	// Load subscriptions
	for _, schemaType := range schema.Types {
		// Skip system types
		if schemaType.BuiltIn {
			continue
		}
		
		// Load subscriptions
		err := r.LoadSubscriptionsFromFile(ctx, schemaType)
		if err != nil {
			return err
		}
	}
	
	return nil
}

func (r *SchemaResolver) LoadSubscriptionsFromFile(ctx context.Context, schemaType *gql.Type) error {
	// Load subscriptions file
	subscriptionsFile, err := os.Open(filepath.Join("subscriptions", strings.ToLower(schemaType.Name) + ".go"))
	if err != nil {
		return err
	}
	defer subscriptionsFile.Close()
	
	// Read subscriptions file
	subscriptionsBytes, err := ioutil.ReadAll(subscriptionsFile)
	if err != nil {
		return err
	}
	
	// Compile subscriptions
	subscriptions, err := r.CompileSubscriptions(ctx, schemaType, string(subscriptionsBytes))
	if err != nil {
		return err
	}
	
	// Set subscriptions
	r.SetSubscriptions(schemaType.Name, subscriptions)
	
	return nil
}

func (r *SchemaResolver) CompileSubscriptions(ctx context.Context, schemaType *gql.Type, subscriptionsCode string) (interface{}, error) {
	// Create resolver context
	resolverContext := &ResolverContext{
		Schema: schemaType,
	}
	
	// Create subscriptions
	subscriptions, err := vm.CreateSubscriptions(resolverContext, subscriptionsCode)
	if err != nil {
		return nil, err
	}
	
	return subscriptions, nil
}

func (r *SchemaResolver) SetSubscriptions(name string, subscriptions interface{}) {
	// Set subscriptions
	r.Subscriptions[name] = subscriptions
}
`))

var serverTemplateInstance = template.Must(template.New("server").Parse(`package main

import (
	"log"
	"net/http"
	"os"
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
`))

var gqlGenYmlTemplate = template.Must(template.New("gqlgenyml").Parse(`schema:
  - graph/*.graphql

exec:
  filename: graph/generated.go
  package: graph

model:
  filename: graph/model/models_gen.go
  package: model

resolver:
  layout: follow-schema
  dir: graph
  package: graph
  filename_template: "{name}.resolvers.go"

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

`))
